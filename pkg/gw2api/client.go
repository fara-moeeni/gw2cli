package gw2api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const BaseURL = "https://api.guildwars2.com/v2"

type Client struct {
	rest   *resty.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	r := resty.New().
		SetBaseURL(BaseURL).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500 || r.StatusCode() == 429
			},
		)

	return &Client{
		rest:   r,
		apiKey: apiKey,
	}
}

func (c *Client) GetSharedInventory() (AccountInventory, error) {
	var inventory AccountInventory
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&inventory).
		Get("/account/inventory")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	return inventory, nil
}

func (c *Client) GetBank() (AccountInventory, error) {
	var bank AccountInventory
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&bank).
		Get("/account/bank")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return bank, nil
}

func (c *Client) GetMaterials() ([]MaterialStorageEntry, error) {
	var materials []MaterialStorageEntry
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&materials).
		Get("/account/materials")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return materials, nil
}

func (c *Client) GetCharacters() ([]Character, error) {
	var characters []Character
	// ids=all fetches full details for all characters in one go
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetQueryParam("ids", "all").
		SetResult(&characters).
		Get("/characters")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return characters, nil
}

func (c *Client) GetItems(ids []int) ([]Item, error) {
	return c.GetItemsWithProgress(ids, nil)
}

func (c *Client) GetItemsWithProgress(ids []int, progress func(int, int)) ([]Item, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	uniqueIDs := make(map[int]bool)
	var list []string
	for _, id := range ids {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			list = append(list, strconv.Itoa(id))
		}
	}

	var allItems []Item
	batchSize := 200

	for i := 0; i < len(list); i += batchSize {
		if progress != nil {
			progress(i, len(list))
		}

		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		batch := list[i:end]
		var batchItems []Item

		resp, err := c.rest.R().
			SetQueryParam("ids", strings.Join(batch, ",")).
			SetResult(&batchItems).
			Get("/items")

		if err != nil {
			return allItems, err
		}
		if resp.IsError() {
			return allItems, fmt.Errorf("API error fetching items: %s", resp.Status())
		}
		allItems = append(allItems, batchItems...)
	}

	if progress != nil {
		progress(len(list), len(list))
	}

	return allItems, nil
}

func (c *Client) GetWallet() ([]WalletCurrency, error) {
	var wallet []WalletCurrency
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&wallet).
		Get("/account/wallet")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return wallet, nil
}

func (c *Client) GetCurrencies(ids []int) ([]Currency, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var list []string
	for _, id := range ids {
		list = append(list, strconv.Itoa(id))
	}

	var currencies []Currency
	resp, err := c.rest.R().
		SetQueryParam("ids", strings.Join(list, ",")).
		SetResult(&currencies).
		Get("/currencies")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return currencies, nil
}

func (c *Client) GetAllItemIDs() ([]int, error) {
	var ids []int
	resp, err := c.rest.R().
		SetResult(&ids).
		Get("/items")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}

func (c *Client) GetCommerceDelivery() (*CommerceDelivery, error) {
	var delivery CommerceDelivery
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&delivery).
		Get("/commerce/delivery")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return &delivery, nil
}

func (c *Client) GetCommercePrices(ids []int) ([]CommercePrice, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var list []string
	for _, id := range ids {
		list = append(list, strconv.Itoa(id))
	}

	var prices []CommercePrice
	resp, err := c.rest.R().
		SetQueryParam("ids", strings.Join(list, ",")).
		SetResult(&prices).
		Get("/commerce/prices")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return prices, nil
}

func (c *Client) GetCommerceTransactions(current bool, buys bool) ([]CommerceTransaction, error) {
	path := "/commerce/transactions/"
	if current {
		path += "current/"
	} else {
		path += "history/"
	}
	if buys {
		path += "buys"
	} else {
		path += "sells"
	}

	var txs []CommerceTransaction
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&txs).
		Get(path)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return txs, nil
}
