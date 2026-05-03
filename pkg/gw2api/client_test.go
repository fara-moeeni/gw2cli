package gw2api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestGetAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		acc := Account{
			Name:         "Tester.1234",
			FractalLevel: 100,
			WvpRank:      500,
		}
		json.NewEncoder(w).Encode(acc)
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.rest.SetBaseURL(server.URL)

	acc, err := client.GetAccount()
	if err != nil {
		t.Fatalf("Failed to get account: %v", err)
	}

	if acc.Name != "Tester.1234" {
		t.Errorf("Expected Tester.1234, got %s", acc.Name)
	}
	if acc.FractalLevel != 100 {
		t.Errorf("Expected 100, got %d", acc.FractalLevel)
	}
}

func TestAuthenticatedRequestSendsBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("missing auth header: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AccountInventory{})
	}))
	defer server.Close()

	client := NewClient("test-key")
	client.rest.SetBaseURL(server.URL)

	if _, err := client.GetBank(); err != nil {
		t.Fatalf("GetBank failed: %v", err)
	}
}

func TestPublicRequestDoesNotRequireBearerToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Fatalf("public request sent auth header: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]string{"behemoth"})
	}))
	defer server.Close()

	client := NewClient("")
	client.rest.SetBaseURL(server.URL)

	ids, err := client.GetWorldBossIDs()
	if err != nil {
		t.Fatalf("GetWorldBossIDs failed: %v", err)
	}
	if len(ids) != 1 || ids[0] != "behemoth" {
		t.Fatalf("unexpected IDs: %#v", ids)
	}
}

func TestGetItemsBatchesRequests(t *testing.T) {
	var batches []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query().Get("ids")
		batches = append(batches, ids)

		w.Header().Set("Content-Type", "application/json")
		var items []Item
		for _, rawID := range strings.Split(ids, ",") {
			id, err := strconv.Atoi(rawID)
			if err != nil {
				t.Fatalf("invalid id %q: %v", rawID, err)
			}
			items = append(items, Item{ID: id, Name: rawID})
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client := NewClient("")
	client.rest.SetBaseURL(server.URL)

	ids := make([]int, 201)
	for i := range ids {
		ids[i] = i + 1
	}

	items, err := client.GetItems(ids)
	if err != nil {
		t.Fatalf("GetItems failed: %v", err)
	}
	if len(items) != 201 {
		t.Fatalf("expected 201 items, got %d", len(items))
	}
	if len(batches) != 2 {
		t.Fatalf("expected 2 batches, got %d", len(batches))
	}
	if got := len(strings.Split(batches[0], ",")); got != 200 {
		t.Fatalf("first batch size = %d, want 200", got)
	}
	if batches[1] != "201" {
		t.Fatalf("unexpected second batch: %q", batches[1])
	}
}
