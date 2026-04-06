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

func (c *Client) GetItemsWithProgress(ids []int, progress func(int, int, []Item)) ([]Item, error) {
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
		
		if progress != nil {
			progress(len(allItems), len(list), batchItems)
		}
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

func (c *Client) GetCoinsToGems(quantity int) (*CommerceExchange, error) {
	var exchange CommerceExchange
	resp, err := c.rest.R().
		SetQueryParam("quantity", strconv.Itoa(quantity)).
		SetResult(&exchange).
		Get("/commerce/exchange/coins")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return &exchange, nil
}

func (c *Client) GetGemsToCoins(quantity int) (*CommerceExchange, error) {
	var exchange CommerceExchange
	resp, err := c.rest.R().
		SetQueryParam("quantity", strconv.Itoa(quantity)).
		SetResult(&exchange).
		Get("/commerce/exchange/gems")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return &exchange, nil
}

func (c *Client) GetLegendaryArmory() ([]LegendaryArmoryItem, error) {
	var items []LegendaryArmoryItem
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&items).
		Get("/account/legendaryarmory")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return items, nil
}

func (c *Client) GetLegendaryArmoryList() ([]int, error) {
	var ids []int
	resp, err := c.rest.R().
		SetResult(&ids).
		Get("/legendaryarmory")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}

func (c *Client) GetAccountRecipes() ([]int, error) {
	var ids []int
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&ids).
		Get("/account/recipes")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}

func (c *Client) GetRecipes(ids []int) ([]Recipe, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var list []string
	for _, id := range ids {
		list = append(list, strconv.Itoa(id))
	}

	var allRecipes []Recipe
	batchSize := 200

	for i := 0; i < len(list); i += batchSize {
		end := i + batchSize
		if end > len(list) {
			end = len(list)
		}

		batch := list[i:end]
		var batchRecipes []Recipe

		resp, err := c.rest.R().
			SetQueryParam("ids", strings.Join(batch, ",")).
			SetResult(&batchRecipes).
			Get("/recipes")

		if err != nil {
			return allRecipes, err
		}
		if resp.IsError() {
			return allRecipes, fmt.Errorf("API error fetching recipes: %s", resp.Status())
		}
		allRecipes = append(allRecipes, batchRecipes...)
	}

	return allRecipes, nil
}

func (c *Client) SearchRecipesByItem(itemID int, input bool) ([]int, error) {
	var ids []int
	req := c.rest.R().SetResult(&ids)
	if input {
		req.SetQueryParam("input", strconv.Itoa(itemID))
	} else {
		req.SetQueryParam("output", strconv.Itoa(itemID))
	}

	resp, err := req.Get("/recipes/search")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}

func (c *Client) GetAccountSkins() ([]int, error) {
	return c.getAccountIDs("/account/skins")
}

func (c *Client) GetAccountDyes() ([]int, error) {
	return c.getAccountIDs("/account/dyes")
}

func (c *Client) GetAccountMinis() ([]int, error) {
	return c.getAccountIDs("/account/minis")
}

func (c *Client) GetAccountMountSkins() ([]int, error) {
	return c.getAccountIDs("/account/mounts/skins")
}

func (c *Client) GetAccountMountTypes() ([]string, error) {
	var types []string
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&types).
		Get("/account/mounts/types")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return types, nil
}

func (c *Client) GetAccountOutfits() ([]int, error) {
	return c.getAccountIDs("/account/outfits")
}

func (c *Client) GetAccountNovelties() ([]int, error) {
	return c.getAccountIDs("/account/novelties")
}

func (c *Client) GetAccountFinishers() ([]int, error) {
	var finishers []struct {
		ID int `json:"id"`
	}
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&finishers).
		Get("/account/finishers")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}

	var ids []int
	for _, f := range finishers {
		ids = append(ids, f.ID)
	}
	return ids, nil
}

func (c *Client) ResolveSkins(ids []int) ([]Skin, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var allSkins []Skin
	batchSize := 200

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		var strIDs []string
		for _, id := range batch {
			strIDs = append(strIDs, strconv.Itoa(id))
		}

		var batchSkins []Skin
		resp, err := c.rest.R().
			SetQueryParam("ids", strings.Join(strIDs, ",")).
			SetResult(&batchSkins).
			Get("/skins")

		if err != nil {
			return allSkins, err
		}
		if resp.IsError() {
			return allSkins, fmt.Errorf("API error: %s", resp.Status())
		}
		allSkins = append(allSkins, batchSkins...)
	}
	return allSkins, nil
}

func (c *Client) ResolveColors(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/colors", ids)
}

func (c *Client) ResolveMinis(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/minis", ids)
}

func (c *Client) ResolveMountSkins(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/mounts/skins", ids)
}

func (c *Client) ResolveOutfits(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/outfits", ids)
}

func (c *Client) ResolveNovelties(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/novelties", ids)
}

func (c *Client) ResolveFinishers(ids []int) ([]NamedEntity, error) {
	return c.resolveEntities("/finishers", ids)
}

func (c *Client) getAccountIDs(path string) ([]int, error) {
	var ids []int
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&ids).
		Get(path)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}

func (c *Client) resolveEntities(path string, ids []int) ([]NamedEntity, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var allItems []NamedEntity
	batchSize := 200

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		var strIDs []string
		for _, id := range batch {
			strIDs = append(strIDs, strconv.Itoa(id))
		}

		var batchItems []NamedEntity
		resp, err := c.rest.R().
			SetQueryParam("ids", strings.Join(strIDs, ",")).
			SetResult(&batchItems).
			Get(path)

		if err != nil {
			return allItems, err
		}
		if resp.IsError() {
			return allItems, fmt.Errorf("API error: %s", resp.Status())
		}
		allItems = append(allItems, batchItems...)
	}
	return allItems, nil
}

func (c *Client) GetDailyAchievements() (*DailyAchievements, error) {
	var dailies DailyAchievements
	resp, err := c.rest.R().
		SetResult(&dailies).
		Get("/achievements/daily")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return &dailies, nil
}

func (c *Client) GetDailyAchievementsTomorrow() (*DailyAchievements, error) {
	var dailies DailyAchievements
	resp, err := c.rest.R().
		SetResult(&dailies).
		Get("/achievements/daily/tomorrow")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return &dailies, nil
}

func (c *Client) GetAchievements(ids []int) ([]Achievement, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var allAchievements []Achievement
	batchSize := 200

	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		var strIDs []string
		for _, id := range batch {
			strIDs = append(strIDs, strconv.Itoa(id))
		}

		var batchAchievements []Achievement
		resp, err := c.rest.R().
			SetQueryParam("ids", strings.Join(strIDs, ",")).
			SetResult(&batchAchievements).
			Get("/achievements")

		if err != nil {
			return allAchievements, err
		}
		if resp.IsError() {
			return allAchievements, fmt.Errorf("API error: %s", resp.Status())
		}
		allAchievements = append(allAchievements, batchAchievements...)
	}
	return allAchievements, nil
}

func (c *Client) GetAccountWorldBosses() ([]string, error) {
	var bosses []string
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&bosses).
		Get("/account/worldbosses")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return bosses, nil
}

func (c *Client) GetAccountDungeons() ([]string, error) {
	var dungeons []string
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&dungeons).
		Get("/account/dungeons")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return dungeons, nil
}

func (c *Client) GetAccountRaids() ([]string, error) {
	var raids []string
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&raids).
		Get("/account/raids")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return raids, nil
}

func (c *Client) GetDungeons() ([]Dungeon, error) {
	var dungeons []Dungeon
	resp, err := c.rest.R().
		SetQueryParam("ids", "all").
		SetResult(&dungeons).
		Get("/dungeons")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return dungeons, nil
}

func (c *Client) GetRaids() ([]Raid, error) {
	var raids []Raid
	resp, err := c.rest.R().
		SetQueryParam("ids", "all").
		SetResult(&raids).
		Get("/raids")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return raids, nil
}

func (c *Client) GetWizardsVaultDaily() ([]WizardsVaultObjective, error) {
	var response WizardsVaultResponse
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&response).
		Get("/account/wizardsvault/daily")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return response.Objectives, nil
}

func (c *Client) GetWizardsVaultWeekly() ([]WizardsVaultObjective, error) {
	var response WizardsVaultResponse
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&response).
		Get("/account/wizardsvault/weekly")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return response.Objectives, nil
}

func (c *Client) GetWizardsVaultSpecial() ([]WizardsVaultObjective, error) {
	var response WizardsVaultResponse
	resp, err := c.rest.R().
		SetHeader("Authorization", "Bearer "+c.apiKey).
		SetResult(&response).
		Get("/account/wizardsvault/special")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return response.Objectives, nil
}

func (c *Client) GetWorldBossIDs() ([]string, error) {
	var ids []string
	resp, err := c.rest.R().Get("/worldbosses")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.Status())
	}
	return ids, nil
}
