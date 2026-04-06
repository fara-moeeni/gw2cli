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

type AchievementCacheEntry struct {
	ID          int                       `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Tiers       []gw2api.AchievementTier  `json:"tiers"`
}

type AchievementCache struct {
	Achievements []AchievementCacheEntry `json:"achievements"`
}

func (s *Service) EnsureAchievementCache(force bool) error {
	cachePath, err := GetAchievementCachePath()
	if err != nil {
		return err
	}

	allIDs, err := s.client.GetAllAchievementIDs()
	if err != nil {
		return err
	}

	var currentCache AchievementCache
	if data, err := os.ReadFile(cachePath); err == nil {
		_ = json.Unmarshal(data, &currentCache)
	}

	if len(currentCache.Achievements) >= len(allIDs) && !force {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	seen := make(map[int]bool)
	for _, a := range currentCache.Achievements {
		seen[a.ID] = true
	}

	var missingIDs []int
	for _, id := range allIDs {
		if !seen[id] {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) == 0 && !force {
		return nil
	}

	fmt.Println("fetching achievement list...")
	fmt.Println("resolving names...")
	_, err = s.client.GetAchievementsWithProgress(missingIDs, func(current, total int, newAch []gw2api.Achievement) {
		for _, a := range newAch {
			currentCache.Achievements = append(currentCache.Achievements, AchievementCacheEntry{
				ID:          a.ID,
				Name:        a.Name,
				Description: a.Description,
				Tiers:       a.Tiers,
			})
		}

		if data, errMarshal := json.Marshal(currentCache); errMarshal == nil {
			_ = os.WriteFile(cachePath, data, 0644)
		}

		pct := float64(len(currentCache.Achievements)) / float64(len(allIDs)) * 100
		barSize := 30
		pos := int(float64(barSize) * (float64(len(currentCache.Achievements)) / float64(len(allIDs))))
		bar := strings.Repeat("=", pos)
		if pos < barSize {
			bar += ">" + strings.Repeat(" ", barSize-pos-1)
		}
		fmt.Printf("\rProgress: [%s] %.1f%% (%d/%d) ", bar, pct, len(currentCache.Achievements), len(allIDs))
	})

	fmt.Printf("\ndone. cached %d achievements to %s\n", len(currentCache.Achievements), cachePath)
	return err
}

func (s *Service) LoadAchievementCache() (*AchievementCache, error) {
	cachePath, err := GetAchievementCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("achievement cache not found, run 'achievements update-cache' to build it")
		}
		return nil, err
	}

	var cache AchievementCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

func (s *Service) CheckAchievementCacheStatus() error {
	cachePath, err := GetAchievementCachePath()
	if err != nil {
		return err
	}

	info, err := os.Stat(cachePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("achievement cache not found, run 'achievements update-cache' to build it")
	}
	if err != nil {
		return err
	}

	if time.Since(info.ModTime()) > 7*24*time.Hour {
		fmt.Println("warning: achievement cache is 7+ days old, run 'achievements update-cache' to refresh")
	}

	return nil
}

func GetAchievementCachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "gw2cli", "achievements.json"), nil
}
