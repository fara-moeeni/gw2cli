package gw2api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

const BaseURL = "https://api.guildwars2.com/v2"

type Client struct {
	rest   *resty.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		rest:   resty.New().SetBaseURL(BaseURL),
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
	if len(ids) == 0 {
		return nil, nil
	}

	// Deduplicate IDs before requesting
	uniqueIDs := make(map[int]bool)
	var list []string
	for _, id := range ids {
		if !uniqueIDs[id] {
			uniqueIDs[id] = true
			list = append(list, strconv.Itoa(id))
		}
	}

	// Chunk requests if there are too many items (API limit is ~200)
	var allItems []Item
	batchSize := 200

	for i := 0; i < len(list); i += batchSize {
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
			return nil, err
		}
		if resp.IsError() {
			return nil, fmt.Errorf("API error fetching items: %s", resp.Status())
		}
		allItems = append(allItems, batchItems...)
	}

	return allItems, nil
}