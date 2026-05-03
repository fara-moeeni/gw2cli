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
	if !force {
		data, err := os.ReadFile(cachePath)
		if err == nil {
			if err := json.Unmarshal(data, &currentCache); err != nil {
				return fmt.Errorf("failed to read item cache: %w", err)
			}
		}
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

	if len(missingIDs) == 0 {
		return nil
	}

	fmt.Println("fetching item list...")
	fmt.Println("resolving names...")
	var writeErr error
	_, err = s.client.GetItemsWithProgress(missingIDs, func(current, total int, newItems []gw2api.Item) {
		for _, item := range newItems {
			currentCache.Items = append(currentCache.Items, CacheEntry{
				ID:   item.ID,
				Name: item.Name,
				Type: item.Type,
			})
		}

		if writeErr == nil {
			writeErr = writeJSONFile(cachePath, currentCache)
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

	if err != nil {
		return fmt.Errorf("caching interrupted: %w", err)
	}
	if writeErr != nil {
		return fmt.Errorf("failed to save item cache: %w", writeErr)
	}

	fmt.Printf("\ndone. cached %d items to %s\n", len(currentCache.Items), cachePath)
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
		return fmt.Errorf("item cache not found, run 'cache update' to build it")
	}
	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > 7*24*time.Hour {
		fmt.Println("warning: item cache is 7+ days old, run 'cache update' to refresh")
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

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	encoder := json.NewEncoder(tmp)
	if err := encoder.Encode(value); err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return fmt.Errorf("%w; failed to close temporary cache file: %v", err, closeErr)
		}
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	return nil
}
