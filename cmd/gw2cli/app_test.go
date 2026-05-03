package main

import (
	"strings"
	"testing"
)

func TestRunAllowsVersionWithoutAPIKey(t *testing.T) {
	if err := run([]string{"version"}, ""); err != nil {
		t.Fatalf("version should not require API key: %v", err)
	}
}

func TestRunRequiresAPIKeyForAuthenticatedCommand(t *testing.T) {
	err := run([]string{"account"}, "")
	if err == nil {
		t.Fatal("account should require an API key")
	}
	if !strings.Contains(err.Error(), "GW2_API_KEY is not set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunValidatesExchangeAmount(t *testing.T) {
	err := run([]string{"exchange", "coins", "not-a-number"}, "")
	if err == nil {
		t.Fatal("exchange should reject invalid amounts")
	}
	if !strings.Contains(err.Error(), "positive integer") {
		t.Fatalf("unexpected error: %v", err)
	}
}
