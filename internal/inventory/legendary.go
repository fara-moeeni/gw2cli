package inventory

import (
	"fmt"
	"strings"
)

type LegendaryItem struct {
	ID    int
	Name  string
	Type  string
	Count int
}

func (s *Service) GetLegendaryArmory(term string) ([]LegendaryItem, error) {
	apiItems, err := s.client.GetLegendaryArmory()
	if err != nil {
		return nil, err
	}

	cache, err := s.LoadCache()
	if err != nil {
		// According to requirements, standard cache error
		return nil, fmt.Errorf("item cache not found, run -update-cache to build it")
	}

	cacheMap := make(map[int]CacheEntry)
	for _, item := range cache.Items {
		cacheMap[item.ID] = item
	}

	term = strings.ToLower(term)
	var results []LegendaryItem
	for _, apiItem := range apiItems {
		name := fmt.Sprintf("Unknown (%d)", apiItem.ID)
		itemType := "Unknown"

		if val, ok := cacheMap[apiItem.ID]; ok {
			name = val.Name
			itemType = val.Type
		}

		if term != "" && !strings.Contains(strings.ToLower(name), term) {
			continue
		}

		results = append(results, LegendaryItem{
			ID:    apiItem.ID,
			Name:  name,
			Type:  itemType,
			Count: apiItem.Count,
		})
	}

	return results, nil
}
