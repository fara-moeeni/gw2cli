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

func (s *Service) EnsureCache() error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(cachePath); err == nil {
		return nil
	}

	fmt.Println("Building local item cache (this may take a minute on first run)...")
	allIDs, err := s.client.GetAllItemIDs()
	if err != nil {
		return err
	}

	items, err := s.client.GetItemsWithProgress(allIDs, func(current, total int) {
		// Only show progress if verbose is enabled
		if !s.Verbose {
			return
		}
		pct := float64(current) / float64(total) * 100
		barSize := 30
		pos := int(float64(barSize) * (float64(current) / float64(total)))
		bar := strings.Repeat("=", pos)
		if pos < barSize {
			bar += ">" + strings.Repeat(" ", barSize-pos-1)
		}
		fmt.Printf("\rDownloading item data: [%s] %.1f%% (%d/%d) ", bar, pct, current, total)
	})
	if s.Verbose {
		fmt.Println()
	}

	if err != nil {
		return err
	}

	var entries []CacheEntry
	for _, item := range items {
		entries = append(entries, CacheEntry{ID: item.ID, Name: item.Name})
	}

	cache := ItemCache{Items: entries}
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
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

	// Parallel search for speed
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
