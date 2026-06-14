package inventory

import (
	"os"
	"path/filepath"
	"testing"

	"gw2cli/pkg/gw2api"
)

func TestSearchFiltersByTermTypeAndLocation(t *testing.T) {
	items := []ItemDetail{
		{
			ID:   1,
			Name: "Mystic Sword",
			Type: "Weapon",
			Locations: []ItemLocation{
				{Source: "Alice", Detail: "Bag 1", Count: 1},
				{Source: "Bank", Detail: "Tab 1", Count: 2},
			},
		},
		{
			ID:        2,
			Name:      "Cloth Coat",
			Type:      "Armor",
			Locations: []ItemLocation{{Source: "Bob", Detail: "Equipped: Coat", Count: 1}},
		},
	}

	results := Search(items, FilterCriteria{
		SearchTerm: "sword",
		Type:       "Weapon",
		Character:  "Alice",
	})

	if len(results) != 1 {
		t.Fatalf("expected one result, got %d", len(results))
	}
	if results[0].Name != "Mystic Sword" {
		t.Fatalf("unexpected item: %s", results[0].Name)
	}
	if len(results[0].Locations) != 1 || results[0].Locations[0].Source != "Alice" {
		t.Fatalf("expected only Alice location, got %#v", results[0].Locations)
	}
}

func TestFilterWalletMatchesNameAndID(t *testing.T) {
	wallet := []WalletEntry{
		{ID: 1, Name: "Coin", Value: 12345},
		{ID: 78, Name: "Fine Rift Essence", Value: 250},
		{ID: 79, Name: "Rare Rift Essence", Value: 12},
	}

	byName := FilterWallet(wallet, "Fine Rift")
	if len(byName) != 1 || byName[0].ID != 78 {
		t.Fatalf("expected Fine Rift Essence by name, got %#v", byName)
	}

	byID := FilterWallet(wallet, "78")
	if len(byID) != 1 || byID[0].Name != "Fine Rift Essence" {
		t.Fatalf("expected Fine Rift Essence by ID, got %#v", byID)
	}
}

func TestFilterWalletHandlesNoMatchAndEmptyTerm(t *testing.T) {
	wallet := []WalletEntry{
		{ID: 78, Name: "Fine Rift Essence", Value: 250},
	}

	if results := FilterWallet(wallet, "magnetite"); len(results) != 0 {
		t.Fatalf("expected no matches, got %#v", results)
	}

	results := FilterWallet(wallet, " ")
	if len(results) != len(wallet) || results[0].ID != wallet[0].ID {
		t.Fatalf("expected empty term to return wallet entries, got %#v", results)
	}
}

func TestMapToProgressComputesTiersAndAP(t *testing.T) {
	service := &Service{}
	progress := service.mapToProgress(
		AchievementCacheEntry{
			ID:   10,
			Name: "Progressive",
			Tiers: []gw2api.AchievementTier{
				{Count: 5, Points: 1},
				{Count: 10, Points: 2},
			},
		},
		gw2api.AccountAchievement{ID: 10, Current: 5},
	)

	if progress.Points != 1 {
		t.Fatalf("expected 1 AP, got %d", progress.Points)
	}
	if progress.TierStatus != "1/2 tiers" {
		t.Fatalf("unexpected tier status: %s", progress.TierStatus)
	}
	if progress.StatusSymbol != "[~]" {
		t.Fatalf("unexpected status: %s", progress.StatusSymbol)
	}
}

func TestEnsureCacheForceRefreshesWithoutDuplicates(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cacheDir := filepath.Join(tmp, ".config", "gw2cli")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}
	cachePath := filepath.Join(cacheDir, "items.json")
	if err := os.WriteFile(cachePath, []byte(`{"items":[{"id":1,"name":"Old","type":"OldType"}]}`), 0644); err != nil {
		t.Fatal(err)
	}

	client := &fakeClient{
		itemIDs: []int{1, 2},
		items: map[int]gw2api.Item{
			1: {ID: 1, Name: "New", Type: "Weapon"},
			2: {ID: 2, Name: "Second", Type: "Armor"},
		},
	}
	service := NewService(client)

	if err := service.EnsureCache(true); err != nil {
		t.Fatalf("EnsureCache failed: %v", err)
	}

	cache, err := service.LoadCache()
	if err != nil {
		t.Fatal(err)
	}
	if len(cache.Items) != 2 {
		t.Fatalf("expected 2 cache entries, got %d: %#v", len(cache.Items), cache.Items)
	}
	if cache.Items[0].Name != "New" {
		t.Fatalf("force refresh did not replace stale item: %#v", cache.Items)
	}
}

func TestEnsureCacheResumesMissingEntries(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	cacheDir := filepath.Join(tmp, ".config", "gw2cli")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		t.Fatal(err)
	}
	cachePath := filepath.Join(cacheDir, "items.json")
	if err := os.WriteFile(cachePath, []byte(`{"items":[{"id":1,"name":"Existing","type":"Weapon"}]}`), 0644); err != nil {
		t.Fatal(err)
	}

	client := &fakeClient{
		itemIDs: []int{1, 2},
		items: map[int]gw2api.Item{
			1: {ID: 1, Name: "Should Not Fetch", Type: "Weapon"},
			2: {ID: 2, Name: "Missing", Type: "Armor"},
		},
	}
	service := NewService(client)

	if err := service.EnsureCache(false); err != nil {
		t.Fatalf("EnsureCache failed: %v", err)
	}

	if len(client.requestedItemIDs) != 1 || client.requestedItemIDs[0] != 2 {
		t.Fatalf("expected only missing ID 2 to be fetched, got %#v", client.requestedItemIDs)
	}
}
