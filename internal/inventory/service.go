package inventory

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"gw2cli/pkg/gw2api"
)

type Service struct {
	client     *gw2api.Client
	Verbose    bool
	BuildCache bool
}

func NewService(client *gw2api.Client) *Service {
	return &Service{client: client}
}

// CharacterSummary contains high-level info for listing characters.
type CharacterSummary struct {
	Name       string
	Race       string
	Gender     string
	Profession string
	Level      int
	Created    time.Time
	Age        time.Duration
}

func (s *Service) GetCharacterList() ([]CharacterSummary, error) {
	chars, err := s.client.GetCharacters()
	if err != nil {
		return nil, err
	}

	var summary []CharacterSummary
	for _, c := range chars {
		createdTime, _ := time.Parse(time.RFC3339, c.Created)
		summary = append(summary, CharacterSummary{
			Name:       c.Name,
			Race:       c.Race,
			Gender:     c.Gender,
			Profession: c.Profession,
			Level:      c.Level,
			Created:    createdTime,
			Age:        time.Duration(c.Age) * time.Second,
		})
	}

	sort.Slice(summary, func(i, j int) bool {
		return summary[i].Name < summary[j].Name
	})

	return summary, nil
}

// FetchAll retrieves all items from Bank, Shared, and Characters.
func (s *Service) FetchAll() ([]ItemDetail, error) {
	var wg sync.WaitGroup
	var shared gw2api.AccountInventory
	var bank gw2api.AccountInventory
	var materials []gw2api.MaterialStorageEntry
	var characters []gw2api.Character
	var errShared, errBank, errMaterials, errChars error

	wg.Add(4)
	go func() { defer wg.Done(); shared, errShared = s.client.GetSharedInventory() }()
	go func() { defer wg.Done(); bank, errBank = s.client.GetBank() }()
	go func() { defer wg.Done(); materials, errMaterials = s.client.GetMaterials() }()
	go func() { defer wg.Done(); characters, errChars = s.client.GetCharacters() }()
	wg.Wait()

	// Return error only if all fetches failed, otherwise partial results are acceptable
	if errShared != nil && errBank != nil && errMaterials != nil && errChars != nil {
		return nil, fmt.Errorf("failed to fetch any data: %v, %v, %v, %v", errShared, errBank, errMaterials, errChars)
	}

	// Map ItemID -> List of Locations
	locMap := make(map[int][]ItemLocation)
	var allIDs []int
	seenIDs := make(map[int]bool)

	add := func(id int, count int, source, detail string) {
		if id == 0 {
			return
		}
		if !seenIDs[id] {
			allIDs = append(allIDs, id)
			seenIDs[id] = true
		}
		locMap[id] = append(locMap[id], ItemLocation{Source: source, Detail: detail, Count: count})
	}

	for i, slot := range shared {
		if slot != nil {
			add(slot.ID, slot.Count, "Shared Inventory", fmt.Sprintf("Slot %d", i+1))
		}
	}

	for i, slot := range bank {
		if slot != nil {
			add(slot.ID, slot.Count, "Bank", fmt.Sprintf("Tab %d", (i/30)+1))
		}
	}

	for _, mat := range materials {
		if mat.Count > 0 {
			add(mat.ID, mat.Count, "Material Storage", "")
		}
	}

	for _, char := range characters {
		for bagIdx, bag := range char.Bags {
			if bag != nil {
				for _, slot := range bag.Inventory {
					if slot != nil {
						add(slot.ID, slot.Count, char.Name, fmt.Sprintf("Bag %d", bagIdx+1))
					}
				}
			}
		}
		for _, equip := range char.Equipment {
			add(equip.ID, 1, char.Name, fmt.Sprintf("Equipped: %s", equip.Slot))
		}
	}

	if len(allIDs) == 0 {
		return nil, nil
	}

	apiItems, err := s.client.GetItems(allIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve item details: %w", err)
	}

	var results []ItemDetail
	for _, apiItem := range apiItems {
		results = append(results, ItemDetail{
			ID:        apiItem.ID,
			Name:      apiItem.Name,
			Type:      apiItem.Type,
			Locations: locMap[apiItem.ID],
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// FilterCriteria defines the filtering parameters for inventory search.
type FilterCriteria struct {
	SearchTerm string
	Type       string
	Character  string
}

// Search filters the items based on the provided criteria.
func Search(items []ItemDetail, criteria FilterCriteria) []ItemDetail {
	var filtered []ItemDetail

	term := strings.ToLower(criteria.SearchTerm)
	targetType := strings.ToLower(criteria.Type)
	targetChar := strings.ToLower(criteria.Character)

	for _, item := range items {
		// Strict Type Filter
		if targetType != "" && strings.ToLower(item.Type) != targetType {
			continue
		}

		// Fuzzy Name/ID Filter
		if term != "" {
			match := strings.Contains(strings.ToLower(item.Name), term) ||
				strings.Contains(fmt.Sprintf("%d", item.ID), term)

			// Fallback: match Type if no strict type was requested
			if targetType == "" && strings.Contains(strings.ToLower(item.Type), term) {
				match = true
			}
			if !match {
				continue
			}
		}

		// Location Filter
		if targetChar != "" {
			var matchingLocs []ItemLocation
			for _, loc := range item.Locations {
				fullLoc := strings.ToLower(loc.Source + " " + loc.Detail)
				if strings.Contains(fullLoc, targetChar) {
					matchingLocs = append(matchingLocs, loc)
				}
			}
			if len(matchingLocs) == 0 {
				continue
			}
			// Clone item with only matching locations
			newItem := item
			newItem.Locations = matchingLocs
			filtered = append(filtered, newItem)
			continue
		}

		filtered = append(filtered, item)
	}
	return filtered
}

func GetUniqueTypes(items []ItemDetail) []string {
	typeMap := make(map[string]bool)
	for _, i := range items {
		typeMap[i.Type] = true
	}

	var types []string
	for t := range typeMap {
		types = append(types, t)
	}
	sort.Strings(types)
	return types
}
