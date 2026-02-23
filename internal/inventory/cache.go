package inventory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gw2cli/pkg/gw2api"
)

type CacheEntry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemCache struct {
	Items []CacheEntry `json:"items"`
}

func (s *Service) EnsureCache(force bool) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	allIDs, err := s.client.GetAllItemIDs()
	if err != nil {
		return err
	}

	var currentCache ItemCache
	if data, err := os.ReadFile(cachePath); err == nil {
		_ = json.Unmarshal(data, &currentCache)
	}

	if len(currentCache.Items) >= len(allIDs) && !force {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Track existing IDs to avoid duplicates during resume
	seen := make(map[int]bool)
	for _, item := range currentCache.Items {
		seen[item.ID] = true
	}

	// Filter IDs we still need to fetch
	var missingIDs []int
	for _, id := range allIDs {
		if !seen[id] {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) == 0 && !force {
		return nil
	}

	fmt.Println("Building local item database...")
	_, err = s.client.GetItemsWithProgress(missingIDs, func(current, total int, newItems []gw2api.Item) {
		for _, item := range newItems {
			currentCache.Items = append(currentCache.Items, CacheEntry{ID: item.ID, Name: item.Name})
		}

		// Save to disk immediately
		if data, errMarshal := json.Marshal(currentCache); errMarshal == nil {
			_ = os.WriteFile(cachePath, data, 0644)
		}

		pct := float64(len(currentCache.Items)) / float64(len(allIDs)) * 100
		barSize := 30
		pos := int(float64(barSize) * (float64(len(currentCache.Items)) / float64(len(allIDs))))
		bar := strings.Repeat("=", pos)
		if pos < barSize {
			bar += ">" + strings.Repeat(" ", barSize-pos-1)
		}
		fmt.Printf("\rProgress: [%s] %.1f%% (%d/%d) ", bar, pct, len(currentCache.Items), len(allIDs))
	})

	fmt.Println()
	if err != nil {
		return fmt.Errorf("caching interrupted: %w", err)
	}

	return nil
}

func (s *Service) SearchCache(term string) ([]int, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache ItemCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	term = strings.ToLower(term)
	var ids []int
	var wg sync.WaitGroup
	var mu sync.Mutex

	numWorkers := 4
	chunkSize := (len(cache.Items) + numWorkers - 1) / numWorkers

	for i := 0; i < len(cache.Items); i += chunkSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunkSize
			if end > len(cache.Items) {
				end = len(cache.Items)
			}
			for _, item := range cache.Items[start:end] {
				if strings.Contains(strings.ToLower(item.Name), term) {
					mu.Lock()
					ids = append(ids, item.ID)
					mu.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()

	return ids, nil
}

func getCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gw2cli", "items.json"), nil
}
