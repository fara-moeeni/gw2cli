package inventory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	exists := false
	if _, err := os.Stat(cachePath); err == nil {
		data, err := os.ReadFile(cachePath)
		if err == nil {
			var cache ItemCache
			if json.Unmarshal(data, &cache) == nil {
				if len(cache.Items) >= len(allIDs) {
					if !force {
						return nil
					}
					exists = true
				}
				if s.Verbose || force {
					fmt.Printf("Cache exists but is incomplete (%d/%d items). Updating...\n", len(cache.Items), len(allIDs))
				}
			}
		}
	}

	if !force && !exists {
		return nil
	}

	// Create directory before starting the long download
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	fmt.Println("Building local item database...")
	items, err := s.client.GetItemsWithProgress(allIDs, func(current, total int) {
		pct := float64(current) / float64(total) * 100
		barSize := 30
		pos := int(float64(barSize) * (float64(current) / float64(total)))
		bar := strings.Repeat("=", pos)
		if pos < barSize {
			bar += ">" + strings.Repeat(" ", barSize-pos-1)
		}
		fmt.Printf("\rProgress: [%s] %.1f%% (%d/%d) ", bar, pct, current, total)
	})

	if len(items) > 0 {
		var entries []CacheEntry
		for _, item := range items {
			entries = append(entries, CacheEntry{ID: item.ID, Name: item.Name})
		}
		cache := ItemCache{Items: entries}
		data, errMarshal := json.Marshal(cache)
		if errMarshal != nil {
			return fmt.Errorf("failed to encode cache: %w", errMarshal)
		}
		if errWrite := os.WriteFile(cachePath, data, 0644); errWrite != nil {
			return fmt.Errorf("failed to write cache file: %w", errWrite)
		}
	}

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

	if _, err := os.Stat(cachePath); err != nil {
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
