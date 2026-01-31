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
	client *gw2api.Client
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
	Age        time.Duration // Playtime
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
	
	// Sort by Name
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
	var characters []gw2api.Character
	var errShared, errBank, errChars error

	// 1. Concurrent Fetching
	wg.Add(3)
	go func() { defer wg.Done(); shared, errShared = s.client.GetSharedInventory() }()
	go func() { defer wg.Done(); bank, errBank = s.client.GetBank() }()
	go func() { defer wg.Done(); characters, errChars = s.client.GetCharacters() }()
	wg.Wait()

	// Handle partial failures (log them but continue if possible?)
	// For now, we return the first error found, or we could just log warnings.
	// Let's return error if everything failed, otherwise proceed with what we have.
	if errShared != nil && errBank != nil && errChars != nil {
		return nil, fmt.Errorf("failed to fetch any data: %v, %v, %v", errShared, errBank, errChars)
	}

	// 2. Aggregation
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

	// Process Shared
	for i, slot := range shared {
		if slot != nil {
			add(slot.ID, slot.Count, "Shared Inventory", fmt.Sprintf("Slot %d", i+1))
		}
	}

	// Process Bank
	for i, slot := range bank {
		if slot != nil {
			add(slot.ID, slot.Count, "Bank", fmt.Sprintf("Tab %d", (i/30)+1))
		}
	}

	// Process Characters
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

	// 3. Resolve Item Details (Names, Types)
	apiItems, err := s.client.GetItems(allIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve item details: %w", err)
	}

	// 4. Build Result
	var results []ItemDetail
	for _, apiItem := range apiItems {
		results = append(results, ItemDetail{
			ID:        apiItem.ID,
			Name:      apiItem.Name,
			Type:      apiItem.Type,
			Locations: locMap[apiItem.ID],
		})
	}

	// Sort by Name by default
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// FilterCriteria defines how we want to filter the inventory.
type FilterCriteria struct {
	SearchTerm string // Fuzzy name/ID
	Type       string // Strict type
	Character  string // Source location (Character name or "Bank")
}

// Search filters the provided list of items based on criteria.
func Search(items []ItemDetail, criteria FilterCriteria) []ItemDetail {
	var filtered []ItemDetail

	term := strings.ToLower(criteria.SearchTerm)
	targetType := strings.ToLower(criteria.Type)
	targetChar := strings.ToLower(criteria.Character)

	for _, item := range items {
		// 1. Strict Type Filter
		if targetType != "" && strings.ToLower(item.Type) != targetType {
			continue
		}

		// 2. Name/ID Fuzzy Filter
		if term != "" {
			match := strings.Contains(strings.ToLower(item.Name), term) ||
				strings.Contains(fmt.Sprintf("%d", item.ID), term)
			
			// Legacy fallback: if no specific type requested, allow term to match Type
			if targetType == "" && strings.Contains(strings.ToLower(item.Type), term) {
				match = true
			}
			if !match {
				continue
			}
		}

		// 3. Location Filter (requires iterating through locations)
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
			// Create a copy of the item with ONLY matching locations
			newItem := item
			newItem.Locations = matchingLocs
			filtered = append(filtered, newItem)
			continue
		}

		// If we made it here, add the item
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