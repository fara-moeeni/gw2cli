package inventory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gw2cli/pkg/gw2api"
)

type CacheEntry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ItemCache struct {
	Items []CacheEntry `json:"items"`
}

func (s *Service) EnsureCache(force bool) error {
	cachePath, err := GetCachePath()
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

	fmt.Println("fetching item list...")
	fmt.Println("resolving names...")
	_, err = s.client.GetItemsWithProgress(missingIDs, func(current, total int, newItems []gw2api.Item) {
		for _, item := range newItems {
			currentCache.Items = append(currentCache.Items, CacheEntry{
				ID:   item.ID,
				Name: item.Name,
				Type: item.Type,
			})
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

	fmt.Printf("\ndone. cached %d items to %s\n", len(currentCache.Items), cachePath)
	if err != nil {
		return fmt.Errorf("caching interrupted: %w", err)
	}

	return nil
}

func (s *Service) SearchCache(term string) ([]int, error) {
	cache, err := s.LoadCache()
	if err != nil {
		return nil, err
	}

	term = strings.ToLower(term)
	var ids []int
	for _, item := range cache.Items {
		if strings.Contains(strings.ToLower(item.Name), term) {
			ids = append(ids, item.ID)
		}
	}
	return ids, nil
}

func (s *Service) LoadCache() (*ItemCache, error) {
	cachePath, err := GetCachePath()
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

	return &cache, nil
}

func (s *Service) CheckCacheStatus() error {
	cachePath, err := GetCachePath()
	if err != nil {
		return err
	}

	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("item cache not found, run -update-cache to build it")
	}
	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > 7*24*time.Hour {
		fmt.Println("warning: item cache is 7+ days old, run -update-cache to refresh")
	}

	return nil
}

func (s *Service) FindInCache(term string) ([]CacheEntry, error) {
	cache, err := s.LoadCache()
	if err != nil {
		return nil, err
	}

	term = strings.ToLower(term)
	var matches []CacheEntry
	for _, item := range cache.Items {
		if strings.Contains(strings.ToLower(item.Name), term) {
			matches = append(matches, item)
		}
	}
	return matches, nil
}

func GetCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gw2cli", "items.json"), nil
}
