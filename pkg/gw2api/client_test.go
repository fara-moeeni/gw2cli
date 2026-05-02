package gw2api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
