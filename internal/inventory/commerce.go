package inventory

import (
	"fmt"
	"gw2cli/pkg/gw2api"
)

type CommerceItem struct {
	ID    int
	Name  string
	Count int
}

type CommerceDelivery struct {
	Coins int
	Items []CommerceItem
}

type CommercePrice struct {
	ID        int
	Name      string
	BuyPrice  int
	SellPrice int
}

type CommerceTransaction struct {
	ID        int
	ItemID    int
	Name      string
	Price     int
	Quantity  int
	Created   string
	Purchased string
}

func (s *Service) GetDelivery() (*CommerceDelivery, error) {
	apiDelivery, err := s.client.GetCommerceDelivery()
	if err != nil {
		return nil, err
	}

	var ids []int
	for _, item := range apiDelivery.Items {
		ids = append(ids, item.ID)
	}

	items, err := s.client.GetItems(ids)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[int]gw2api.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	delivery := &CommerceDelivery{
		Coins: apiDelivery.Coins,
	}

	for _, item := range apiDelivery.Items {
		name := fmt.Sprintf("Unknown (%d)", item.ID)
		if val, ok := itemMap[item.ID]; ok {
			name = val.Name
		}
		delivery.Items = append(delivery.Items, CommerceItem{
			ID:    item.ID,
			Name:  name,
			Count: item.Count,
		})
	}

	return delivery, nil
}

func (s *Service) GetPrices(term string) ([]CommercePrice, error) {
	var ids []int
	var termID int
	if _, err := fmt.Sscanf(term, "%d", &termID); err == nil {
		ids = []int{termID}
	} else {
		if s.SkipCache {
			return nil, fmt.Errorf("name search requires local cache. Use an Item ID instead, or remove -no-cache")
		}
		// Use local cache for name search
		if err := s.EnsureCache(); err != nil {
			return nil, fmt.Errorf("failed to ensure cache: %w", err)
		}
		cachedIDs, err := s.SearchCache(term)
		if err != nil {
			return nil, fmt.Errorf("failed to search cache: %w", err)
		}
		ids = cachedIDs
	}

	if len(ids) == 0 {
		return nil, nil
	}

	// Limit to top results for safety if name matched many items
	if len(ids) > 10 {
		ids = ids[:10]
	}

	apiPrices, err := s.client.GetCommercePrices(ids)

	if err != nil {
		return nil, err
	}

	items, err := s.client.GetItems(ids)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[int]gw2api.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	var results []CommercePrice
	for _, p := range apiPrices {
		name := fmt.Sprintf("Unknown (%d)", p.ID)
		if val, ok := itemMap[p.ID]; ok {
			name = val.Name
		}
		results = append(results, CommercePrice{
			ID:        p.ID,
			Name:      name,
			BuyPrice:  p.Buys.UnitPrice,
			SellPrice: p.Sells.UnitPrice,
		})
	}

	return results, nil
}

func (s *Service) GetTransactions(current bool) ([]CommerceTransaction, []CommerceTransaction, error) {
	buys, err := s.client.GetCommerceTransactions(current, true)
	if err != nil {
		return nil, nil, err
	}
	sells, err := s.client.GetCommerceTransactions(current, false)
	if err != nil {
		return nil, nil, err
	}

	var ids []int
	for _, tx := range buys {
		ids = append(ids, tx.ItemID)
	}
	for _, tx := range sells {
		ids = append(ids, tx.ItemID)
	}

	items, err := s.client.GetItems(ids)
	if err != nil {
		return nil, nil, err
	}

	itemMap := make(map[int]gw2api.Item)
	for _, item := range items {
		itemMap[item.ID] = item
	}

	mapTx := func(apiTxs []gw2api.CommerceTransaction) []CommerceTransaction {
		var results []CommerceTransaction
		for _, tx := range apiTxs {
			name := fmt.Sprintf("Unknown (%d)", tx.ItemID)
			if val, ok := itemMap[tx.ItemID]; ok {
				name = val.Name
			}
			results = append(results, CommerceTransaction{
				ID:        tx.ID,
				ItemID:    tx.ItemID,
				Name:      name,
				Price:     tx.Price,
				Quantity:  tx.Quantity,
				Created:   tx.Created,
				Purchased: tx.Purchased,
			})
		}
		return results
	}

	return mapTx(buys), mapTx(sells), nil
}
