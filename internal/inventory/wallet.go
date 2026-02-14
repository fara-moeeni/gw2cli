package inventory

import (
	"fmt"
	"gw2cli/pkg/gw2api"
)

type WalletEntry struct {
	ID          int
	Name        string
	Value       int
	Description string
}

func (s *Service) GetWallet() ([]WalletEntry, error) {
	apiWallet, err := s.client.GetWallet()
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, w := range apiWallet {
		ids = append(ids, w.ID)
	}

	currencies, err := s.client.GetCurrencies(ids)
	if err != nil {
		return nil, err
	}

	currMap := make(map[int]gw2api.Currency)
	for _, c := range currencies {
		currMap[c.ID] = c
	}

	var results []WalletEntry
	for _, w := range apiWallet {
		c, ok := currMap[w.ID]
		name := fmt.Sprintf("Unknown (%d)", w.ID)
		desc := ""
		if ok {
			name = c.Name
			desc = c.Description
		}
		results = append(results, WalletEntry{
			ID:          w.ID,
			Name:        name,
			Value:       w.Value,
			Description: desc,
		})
	}

	return results, nil
}
