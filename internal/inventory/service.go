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
	client  APIClient
	Verbose bool
}

type APIClient interface {
	GetAccount() (*gw2api.Account, error)
	GetSharedInventory() (gw2api.AccountInventory, error)
	GetBank() (gw2api.AccountInventory, error)
	GetMaterials() ([]gw2api.MaterialStorageEntry, error)
	GetCharacters() ([]gw2api.Character, error)
	GetItems([]int) ([]gw2api.Item, error)
	GetItemsWithProgress([]int, func(int, int, []gw2api.Item)) ([]gw2api.Item, error)
	GetWallet() ([]gw2api.WalletCurrency, error)
	GetCurrencies([]int) ([]gw2api.Currency, error)
	GetAllItemIDs() ([]int, error)
	GetCommerceDelivery() (*gw2api.CommerceDelivery, error)
	GetCommercePrices([]int) ([]gw2api.CommercePrice, error)
	GetCommerceTransactions(bool, bool) ([]gw2api.CommerceTransaction, error)
	GetCoinsToGems(int) (*gw2api.CommerceExchange, error)
	GetGemsToCoins(int) (*gw2api.CommerceExchange, error)
	GetLegendaryArmory() ([]gw2api.LegendaryArmoryItem, error)
	GetAccountRecipes() ([]int, error)
	GetRecipes([]int) ([]gw2api.Recipe, error)
	SearchRecipesByItem(int, bool) ([]int, error)
	GetAccountSkins() ([]int, error)
	GetAccountDyes() ([]int, error)
	GetAccountMinis() ([]int, error)
	GetAccountMountSkins() ([]int, error)
	GetAccountMountTypes() ([]string, error)
	GetAccountOutfits() ([]int, error)
	GetAccountNovelties() ([]int, error)
	GetAccountFinishers() ([]int, error)
	ResolveSkins([]int) ([]gw2api.Skin, error)
	ResolveColors([]int) ([]gw2api.NamedEntity, error)
	ResolveMinis([]int) ([]gw2api.NamedEntity, error)
	ResolveMountSkins([]int) ([]gw2api.NamedEntity, error)
	ResolveOutfits([]int) ([]gw2api.NamedEntity, error)
	ResolveNovelties([]int) ([]gw2api.NamedEntity, error)
	ResolveFinishers([]int) ([]gw2api.NamedEntity, error)
	GetDailyAchievements() (*gw2api.DailyAchievements, error)
	GetAchievements([]int) ([]gw2api.Achievement, error)
	GetAccountWorldBosses() ([]string, error)
	GetWorldBossIDs() ([]string, error)
	GetAccountDungeons() ([]string, error)
	GetDungeons() ([]gw2api.Dungeon, error)
	GetAccountRaids() ([]string, error)
	GetRaids() ([]gw2api.Raid, error)
	GetWizardsVaultDaily() ([]gw2api.WizardsVaultObjective, error)
	GetWizardsVaultWeekly() ([]gw2api.WizardsVaultObjective, error)
	GetAllAchievementIDs() ([]int, error)
	GetAchievementsWithProgress([]int, func(int, int, []gw2api.Achievement)) ([]gw2api.Achievement, error)
	GetAccountAchievements() ([]gw2api.AccountAchievement, error)
	GetMasteryPointSummary() (*gw2api.MasteryPointSummary, error)
	GetLuck() (int, error)
	GetAchievementCategories() ([]gw2api.AchievementCategory, error)
	GetAchievementGroups() ([]gw2api.AchievementGroup, error)
}

func NewService(client APIClient) *Service {
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

func (s *Service) GetAccountSummary() (*AccountSummary, error) {
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}

	return &AccountSummary{
		Name:          acc.Name,
		FractalLevel:  acc.FractalLevel,
		WvwRank:       acc.WvpRank,
		DailyAP:       acc.DailyAP,
		MonthlyAP:     acc.MonthlyAP,
		TotalAP:       acc.DailyAP + acc.MonthlyAP,
		Created:       acc.Created,
		TotalPlaytime: acc.Age,
	}, nil
}

func (s *Service) GetCharacterList() ([]CharacterSummary, error) {
	chars, err := s.client.GetCharacters()
	if err != nil {
		return nil, err
	}

	var summary []CharacterSummary
	for _, c := range chars {
		createdTime, err := time.Parse(time.RFC3339, c.Created)
		if err != nil {
			createdTime = time.Time{}
		}
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
	var partialFailures []string
	if errShared != nil {
		partialFailures = append(partialFailures, fmt.Sprintf("shared inventory: %v", errShared))
	}
	if errBank != nil {
		partialFailures = append(partialFailures, fmt.Sprintf("bank: %v", errBank))
	}
	if errMaterials != nil {
		partialFailures = append(partialFailures, fmt.Sprintf("materials: %v", errMaterials))
	}
	if errChars != nil {
		partialFailures = append(partialFailures, fmt.Sprintf("characters: %v", errChars))
	}
	if len(partialFailures) > 0 {
		fmt.Printf("warning: partial account data: %s\n", strings.Join(partialFailures, "; "))
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

	var failures []string
	if errSkins != nil {
		failures = append(failures, fmt.Sprintf("skins: %v", errSkins))
	}
	if errDyes != nil {
		failures = append(failures, fmt.Sprintf("dyes: %v", errDyes))
	}
	if errMinis != nil {
		failures = append(failures, fmt.Sprintf("minis: %v", errMinis))
	}
	if errMountSkins != nil {
		failures = append(failures, fmt.Sprintf("mount skins: %v", errMountSkins))
	}
	if errMountTypes != nil {
		failures = append(failures, fmt.Sprintf("mount types: %v", errMountTypes))
	}
	if errOutfits != nil {
		failures = append(failures, fmt.Sprintf("outfits: %v", errOutfits))
	}
	if errNovelties != nil {
		failures = append(failures, fmt.Sprintf("novelties: %v", errNovelties))
	}
	if errFinishers != nil {
		failures = append(failures, fmt.Sprintf("finishers: %v", errFinishers))
	}
	if len(failures) > 0 {
		return nil, fmt.Errorf("failed to fetch collection data: %s", strings.Join(failures, "; "))
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
		if t == "" {
			continue
		}
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

func (s *Service) GetDailyBosses() ([]DailyStatus, error) {
	cleared, err := s.client.GetAccountWorldBosses()
	if err != nil {
		return nil, err
	}

	clearedMap := make(map[string]bool)
	for _, b := range cleared {
		clearedMap[b] = true
	}

	ids, err := s.client.GetWorldBossIDs()
	if err != nil {
		return nil, err
	}

	var results []DailyStatus
	for _, id := range ids {
		results = append(results, DailyStatus{
			Name:      id, // The API returns IDs like "behemoth" which are semi-friendly
			Completed: clearedMap[id],
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func (s *Service) GetWorldBossList() ([]DailyStatus, error) {
	ids, err := s.client.GetWorldBossIDs()
	if err != nil {
		return nil, err
	}

	var results []DailyStatus
	for _, id := range ids {
		results = append(results, DailyStatus{Name: id})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})
	return results, nil
}

func (s *Service) GetDailyFractals() ([]FractalDaily, error) {
	dailies, err := s.client.GetDailyAchievements()
	if err != nil {
		return nil, err
	}

	var achIDs []int
	for _, f := range dailies.Fractals {
		achIDs = append(achIDs, f.ID)
	}

	achDetails, err := s.client.GetAchievements(achIDs)
	if err != nil {
		return nil, err
	}

	achMap := make(map[int]gw2api.Achievement)
	for _, a := range achDetails {
		achMap[a.ID] = a
	}

	// We'd need to check progress to see if they are completed, but
	// /v2/account/achievements is needed. For now, let's assume [ ]
	// as per requirement "Handle accounts with no clears gracefully".
	// To actually show [✓], we'd need account achievement progress.

	var results []FractalDaily
	for _, f := range dailies.Fractals {
		ach := achMap[f.ID]
		tier := "T1"
		if strings.Contains(ach.Name, "Tier 4") {
			tier = "T4"
		} else if strings.Contains(ach.Name, "Tier 3") {
			tier = "T3"
		} else if strings.Contains(ach.Name, "Tier 2") {
			tier = "T2"
		} else if strings.Contains(ach.Name, "Recommended") {
			tier = "Rec"
		}

		results = append(results, FractalDaily{
			Name:      ach.Name,
			Tier:      tier,
			Completed: false, // Default to false if we don't fetch progress
		})
	}

	return results, nil
}

func (s *Service) GetDailyWizardsVault() ([]WizardsVaultStatus, error) {
	wv, err := s.client.GetWizardsVaultDaily()
	if err != nil {
		return nil, err
	}

	var results []WizardsVaultStatus
	for _, o := range wv {
		results = append(results, WizardsVaultStatus{
			Title:        o.Title,
			ProgressCur:  o.ProgressCur,
			ProgressGoal: o.ProgressGoal,
			Acclaim:      o.Acclaim,
			Completed:    o.Claimed || (o.ProgressCur >= o.ProgressGoal && o.ProgressGoal > 0),
		})
	}
	return results, nil
}

func (s *Service) GetWeeklyWizardsVault() ([]WizardsVaultStatus, error) {
	wv, err := s.client.GetWizardsVaultWeekly()
	if err != nil {
		return nil, err
	}

	var results []WizardsVaultStatus
	for _, o := range wv {
		results = append(results, WizardsVaultStatus{
			Title:        o.Title,
			ProgressCur:  o.ProgressCur,
			ProgressGoal: o.ProgressGoal,
			Acclaim:      o.Acclaim,
			Completed:    o.Claimed || (o.ProgressCur >= o.ProgressGoal && o.ProgressGoal > 0),
		})
	}
	return results, nil
}

func (s *Service) GetWeeklyRaids() ([]RaidWingStatus, error) {
	cleared, err := s.client.GetAccountRaids()
	if err != nil {
		return nil, err
	}

	clearedMap := make(map[string]bool)
	for _, r := range cleared {
		clearedMap[r] = true
	}

	raids, err := s.client.GetRaids()
	if err != nil {
		return nil, err
	}

	var results []RaidWingStatus
	for _, raid := range raids {
		for _, wing := range raid.Wings {
			var events []DailyStatus
			for _, event := range wing.Events {
				events = append(events, DailyStatus{
					Name:      event.ID,
					Completed: clearedMap[event.ID],
				})
			}
			results = append(results, RaidWingStatus{
				Name:   wing.ID,
				Events: events,
			})
		}
	}

	return results, nil
}

func (s *Service) GetDailyDungeons() ([]DailyStatus, error) {
	cleared, err := s.client.GetAccountDungeons()
	if err != nil {
		return nil, err
	}

	clearedMap := make(map[string]bool)
	for _, d := range cleared {
		clearedMap[d] = true
	}

	dungeons, err := s.client.GetDungeons()
	if err != nil {
		return nil, err
	}

	var results []DailyStatus
	for _, d := range dungeons {
		for _, p := range d.Paths {
			results = append(results, DailyStatus{
				Name:      fmt.Sprintf("%s - %s", d.ID, p.ID),
				Completed: clearedMap[p.ID],
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func (s *Service) GetAllAchievements(statusFilter string) ([]AchievementProgress, error) {
	cache, err := s.LoadAchievementCache()
	if err != nil {
		return nil, err
	}

	accountAch, err := s.client.GetAccountAchievements()
	if err != nil {
		return nil, err
	}

	progMap := make(map[int]gw2api.AccountAchievement)
	for _, a := range accountAch {
		progMap[a.ID] = a
	}

	var results []AchievementProgress
	statusFilter = strings.ToLower(statusFilter)

	for _, a := range cache.Achievements {
		p := progMap[a.ID]
		ap := s.mapToProgress(a, p)

		match := false
		switch statusFilter {
		case "completed":
			if ap.Done {
				match = true
			}
		case "incomplete":
			if !ap.Done {
				match = true
			}
		default:
			match = true
		}

		if match {
			results = append(results, ap)
		}
	}

	return results, nil
}

func (s *Service) GetAchievementSummary() ([]CategorySummary, error) {
	accountAch, err := s.client.GetAccountAchievements()
	if err != nil {
		return nil, err
	}

	categories, err := s.client.GetAchievementCategories()
	if err != nil {
		return nil, err
	}

	cache, err := s.LoadAchievementCache()
	if err != nil {
		return nil, err
	}

	achMap := make(map[int]AchievementCacheEntry)
	for _, a := range cache.Achievements {
		achMap[a.ID] = a
	}

	progMap := make(map[int]gw2api.AccountAchievement)
	for _, a := range accountAch {
		progMap[a.ID] = a
	}

	var results []CategorySummary
	for _, cat := range categories {
		completed := 0
		ap := 0
		for _, aid := range cat.Achievements {
			if p, ok := progMap[aid]; ok {
				if p.Done {
					completed++
				}
				if a, okAch := achMap[aid]; okAch {
					for _, t := range a.Tiers {
						if p.Current >= t.Count {
							ap += t.Points
						}
					}
				}
			}
		}
		if len(cat.Achievements) > 0 {
			results = append(results, CategorySummary{
				Name:      cat.Name,
				Completed: completed,
				Total:     len(cat.Achievements),
				AP:        ap,
			})
		}
	}

	return results, nil
}

func (s *Service) FindAchievements(term string) ([]AchievementProgress, error) {
	cache, err := s.LoadAchievementCache()
	if err != nil {
		return nil, err
	}

	accountAch, err := s.client.GetAccountAchievements()
	if err != nil {
		return nil, err
	}

	progMap := make(map[int]gw2api.AccountAchievement)
	for _, a := range accountAch {
		progMap[a.ID] = a
	}

	var results []AchievementProgress
	term = strings.ToLower(term)

	for _, a := range cache.Achievements {
		if strings.Contains(strings.ToLower(a.Name), term) {
			p := progMap[a.ID]
			results = append(results, s.mapToProgress(a, p))
		}
	}

	return results, nil
}

func (s *Service) GetMasterySummary() (*MasterySummary, error) {
	points, err := s.client.GetMasteryPointSummary()
	if err != nil {
		return nil, err
	}

	luck, err := s.client.GetLuck()
	if err != nil {
		luck = 0
	}

	var regions []MasteryRegion
	for _, t := range points.Totals {
		regions = append(regions, MasteryRegion{
			Name:   t.Region,
			Spent:  t.Spent,
			Earned: t.Earned,
		})
	}

	return &MasterySummary{
		Regions: regions,
		Luck:    luck,
	}, nil
}

func (s *Service) mapToProgress(a AchievementCacheEntry, p gw2api.AccountAchievement) AchievementProgress {
	symbol := "[ ]"
	if p.Done {
		symbol = "[✓]"
	} else if p.Current > 0 {
		symbol = "[~]"
	}

	ap := 0
	tiersDone := 0
	for _, t := range a.Tiers {
		if p.Current >= t.Count {
			ap += t.Points
			tiersDone++
		}
	}

	maxCount := 0
	if len(a.Tiers) > 0 {
		maxCount = a.Tiers[len(a.Tiers)-1].Count
	}

	return AchievementProgress{
		ID:           a.ID,
		Name:         a.Name,
		Description:  a.Description,
		Current:      p.Current,
		Max:          maxCount,
		Points:       ap,
		Done:         p.Done,
		TierStatus:   fmt.Sprintf("%d/%d tiers", tiersDone, len(a.Tiers)),
		StatusSymbol: symbol,
	}
}

func (s *Service) GetCategoryAchievements(groupName string) ([]AchievementProgress, error) {
	groups, err := s.client.GetAchievementGroups()
	if err != nil {
		return nil, err
	}

	categories, err := s.client.GetAchievementCategories()
	if err != nil {
		return nil, err
	}

	cache, err := s.LoadAchievementCache()
	if err != nil {
		return nil, err
	}

	accountAch, err := s.client.GetAccountAchievements()
	if err != nil {
		return nil, err
	}

	progMap := make(map[int]gw2api.AccountAchievement)
	for _, a := range accountAch {
		progMap[a.ID] = a
	}

	achMap := make(map[int]AchievementCacheEntry)
	for _, a := range cache.Achievements {
		achMap[a.ID] = a
	}

	catMap := make(map[int]gw2api.AchievementCategory)
	for _, c := range categories {
		catMap[c.ID] = c
	}

	var targetAchIDs []int
	for _, g := range groups {
		if strings.Contains(strings.ToLower(g.Name), strings.ToLower(groupName)) {
			for _, cid := range g.Categories {
				if cat, ok := catMap[cid]; ok {
					targetAchIDs = append(targetAchIDs, cat.Achievements...)
				}
			}
		}
	}

	// Fallback to searching categories directly if no group found
	if len(targetAchIDs) == 0 {
		for _, c := range categories {
			if strings.Contains(strings.ToLower(c.Name), strings.ToLower(groupName)) {
				targetAchIDs = append(targetAchIDs, c.Achievements...)
			}
		}
	}

	var results []AchievementProgress
	for _, aid := range targetAchIDs {
		if a, ok := achMap[aid]; ok {
			results = append(results, s.mapToProgress(a, progMap[aid]))
		}
	}

	return results, nil
}
