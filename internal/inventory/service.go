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

func (s *Service) GetUnlockedRecipes(term string) ([]RecipeDetail, error) {
	unlockedIDs, err := s.client.GetAccountRecipes()
	if err != nil {
		return nil, err
	}

	if len(unlockedIDs) == 0 {
		return nil, nil
	}

	recipes, err := s.client.GetRecipes(unlockedIDs)
	if err != nil {
		return nil, err
	}

	cache, err := s.LoadCache()
	if err != nil {
		return nil, fmt.Errorf("item cache not found, run 'cache update' to build it")
	}

	cacheMap := make(map[int]string)
	for _, item := range cache.Items {
		cacheMap[item.ID] = item.Name
	}

	var results []RecipeDetail
	term = strings.ToLower(term)

	for _, r := range recipes {
		name := cacheMap[r.OutputItemID]
		if name == "" {
			name = fmt.Sprintf("Unknown (%d)", r.OutputItemID)
		}

		if term != "" && !strings.Contains(strings.ToLower(name), term) {
			continue
		}

		results = append(results, RecipeDetail{
			ID:         r.ID,
			OutputName: name,
			Discipline: strings.Join(r.Disciplines, ", "),
			Rating:     r.MinRating,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].OutputName < results[j].OutputName
	})

	return results, nil
}

func (s *Service) SearchRecipesByIngredient(term string) ([]RecipeDetail, error) {
	matchingItems, err := s.FindInCache(term)
	if err != nil {
		return nil, err
	}

	unlockedIDs, err := s.client.GetAccountRecipes()
	if err != nil {
		return nil, err
	}

	unlockedMap := make(map[int]bool)
	for _, id := range unlockedIDs {
		unlockedMap[id] = true
	}

	var allRecipeIDs []int
	for _, item := range matchingItems {
		ids, err := s.client.SearchRecipesByItem(item.ID, true)
		if err == nil {
			for _, id := range ids {
				if unlockedMap[id] {
					allRecipeIDs = append(allRecipeIDs, id)
				}
			}
		}
	}

	if len(allRecipeIDs) == 0 {
		return nil, nil
	}

	recipes, err := s.client.GetRecipes(allRecipeIDs)
	if err != nil {
		return nil, err
	}

	cache, err := s.LoadCache()
	if err != nil {
		return nil, fmt.Errorf("item cache not found, run 'cache update' to build it")
	}

	cacheMap := make(map[int]string)
	for _, item := range cache.Items {
		cacheMap[item.ID] = item.Name
	}

	var results []RecipeDetail
	for _, r := range recipes {
		name := cacheMap[r.OutputItemID]
		if name == "" {
			name = fmt.Sprintf("Unknown (%d)", r.OutputItemID)
		}

		var ingredients []RecipeIngredientDetail
		for _, ing := range r.Ingredients {
			ingName := cacheMap[ing.ItemID]
			if ingName == "" {
				ingName = fmt.Sprintf("Unknown (%d)", ing.ItemID)
			}
			ingredients = append(ingredients, RecipeIngredientDetail{
				ItemID: ing.ItemID,
				Name:   ingName,
				Count:  ing.Count,
			})
		}

		results = append(results, RecipeDetail{
			ID:          r.ID,
			OutputName:  name,
			Discipline:  strings.Join(r.Disciplines, ", "),
			Rating:      r.MinRating,
			Ingredients: ingredients,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].OutputName < results[j].OutputName
	})

	return results, nil
}

func (s *Service) GetCollectionSummary() (*CollectionSummary, error) {
	var wg sync.WaitGroup
	var skins, dyes, minis, mountSkins, outfits, novelties, finishers []int
	var mountTypes []string
	var errSkins, errDyes, errMinis, errMountSkins, errMountTypes, errOutfits, errNovelties, errFinishers error

	wg.Add(8)
	go func() { defer wg.Done(); skins, errSkins = s.client.GetAccountSkins() }()
	go func() { defer wg.Done(); dyes, errDyes = s.client.GetAccountDyes() }()
	go func() { defer wg.Done(); minis, errMinis = s.client.GetAccountMinis() }()
	go func() { defer wg.Done(); mountSkins, errMountSkins = s.client.GetAccountMountSkins() }()
	go func() { defer wg.Done(); mountTypes, errMountTypes = s.client.GetAccountMountTypes() }()
	go func() { defer wg.Done(); outfits, errOutfits = s.client.GetAccountOutfits() }()
	go func() { defer wg.Done(); novelties, errNovelties = s.client.GetAccountNovelties() }()
	go func() { defer wg.Done(); finishers, errFinishers = s.client.GetAccountFinishers() }()
	wg.Wait()

	if errSkins != nil || errDyes != nil || errMinis != nil || errMountSkins != nil || errMountTypes != nil || errOutfits != nil || errNovelties != nil || errFinishers != nil {
		return nil, fmt.Errorf("failed to fetch collection data")
	}

	return &CollectionSummary{
		Skins:     len(skins),
		Dyes:      len(dyes),
		Minis:     len(minis),
		Mounts:    len(mountSkins) + len(mountTypes),
		Outfits:   len(outfits),
		Novelties: len(novelties),
		Finishers: len(finishers),
	}, nil
}

func (s *Service) GetCollectionSkins(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountSkins()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveSkins(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name, Type: r.Type})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionDyes(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountDyes()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveColors(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionMinis(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountMinis()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveMinis(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionMounts(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountMountSkins()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveMountSkins(ids)
	if err != nil {
		return nil, err
	}
	types, err := s.client.GetAccountMountTypes()
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name, Type: "Skin"})
	}
	for _, t := range types {
		// Capitalize mount type
		name := strings.ToUpper(t[:1]) + t[1:]
		items = append(items, CollectionItem{Name: name, Type: "Type"})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionOutfits(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountOutfits()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveOutfits(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionNovelties(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountNovelties()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveNovelties(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) GetCollectionFinishers(term string) ([]CollectionItem, error) {
	ids, err := s.client.GetAccountFinishers()
	if err != nil {
		return nil, err
	}
	resolved, err := s.client.ResolveFinishers(ids)
	if err != nil {
		return nil, err
	}
	var items []CollectionItem
	for _, r := range resolved {
		items = append(items, CollectionItem{Name: r.Name})
	}
	return s.filterAndSort(items, term), nil
}

func (s *Service) filterAndSort(items []CollectionItem, term string) []CollectionItem {
	var filtered []CollectionItem
	term = strings.ToLower(term)
	for _, item := range items {
		if term == "" || strings.Contains(strings.ToLower(item.Name), term) {
			filtered = append(filtered, item)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Name < filtered[j].Name
	})
	return filtered
}
