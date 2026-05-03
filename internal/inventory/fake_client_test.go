package inventory

import (
	"fmt"

	"gw2cli/pkg/gw2api"
)

type fakeClient struct {
	itemIDs          []int
	items            map[int]gw2api.Item
	requestedItemIDs []int
}

func (f *fakeClient) GetAccount() (*gw2api.Account, error) {
	return nil, errNotImplemented("GetAccount")
}
func (f *fakeClient) GetSharedInventory() (gw2api.AccountInventory, error) {
	return nil, errNotImplemented("GetSharedInventory")
}
func (f *fakeClient) GetBank() (gw2api.AccountInventory, error) {
	return nil, errNotImplemented("GetBank")
}
func (f *fakeClient) GetMaterials() ([]gw2api.MaterialStorageEntry, error) {
	return nil, errNotImplemented("GetMaterials")
}
func (f *fakeClient) GetCharacters() ([]gw2api.Character, error) {
	return nil, errNotImplemented("GetCharacters")
}
func (f *fakeClient) GetItems(ids []int) ([]gw2api.Item, error) {
	return f.resolveItems(ids), nil
}
func (f *fakeClient) GetItemsWithProgress(ids []int, progress func(int, int, []gw2api.Item)) ([]gw2api.Item, error) {
	f.requestedItemIDs = append(f.requestedItemIDs, ids...)
	items := f.resolveItems(ids)
	if progress != nil {
		progress(len(items), len(ids), items)
	}
	return items, nil
}
func (f *fakeClient) GetWallet() ([]gw2api.WalletCurrency, error) {
	return nil, errNotImplemented("GetWallet")
}
func (f *fakeClient) GetCurrencies([]int) ([]gw2api.Currency, error) {
	return nil, errNotImplemented("GetCurrencies")
}
func (f *fakeClient) GetAllItemIDs() ([]int, error) { return f.itemIDs, nil }
func (f *fakeClient) GetCommerceDelivery() (*gw2api.CommerceDelivery, error) {
	return nil, errNotImplemented("GetCommerceDelivery")
}
func (f *fakeClient) GetCommercePrices([]int) ([]gw2api.CommercePrice, error) {
	return nil, errNotImplemented("GetCommercePrices")
}
func (f *fakeClient) GetCommerceTransactions(bool, bool) ([]gw2api.CommerceTransaction, error) {
	return nil, errNotImplemented("GetCommerceTransactions")
}
func (f *fakeClient) GetCoinsToGems(int) (*gw2api.CommerceExchange, error) {
	return nil, errNotImplemented("GetCoinsToGems")
}
func (f *fakeClient) GetGemsToCoins(int) (*gw2api.CommerceExchange, error) {
	return nil, errNotImplemented("GetGemsToCoins")
}
func (f *fakeClient) GetLegendaryArmory() ([]gw2api.LegendaryArmoryItem, error) {
	return nil, errNotImplemented("GetLegendaryArmory")
}
func (f *fakeClient) GetAccountRecipes() ([]int, error) {
	return nil, errNotImplemented("GetAccountRecipes")
}
func (f *fakeClient) GetRecipes([]int) ([]gw2api.Recipe, error) {
	return nil, errNotImplemented("GetRecipes")
}
func (f *fakeClient) SearchRecipesByItem(int, bool) ([]int, error) {
	return nil, errNotImplemented("SearchRecipesByItem")
}
func (f *fakeClient) GetAccountSkins() ([]int, error) {
	return nil, errNotImplemented("GetAccountSkins")
}
func (f *fakeClient) GetAccountDyes() ([]int, error) {
	return nil, errNotImplemented("GetAccountDyes")
}
func (f *fakeClient) GetAccountMinis() ([]int, error) {
	return nil, errNotImplemented("GetAccountMinis")
}
func (f *fakeClient) GetAccountMountSkins() ([]int, error) {
	return nil, errNotImplemented("GetAccountMountSkins")
}
func (f *fakeClient) GetAccountMountTypes() ([]string, error) {
	return nil, errNotImplemented("GetAccountMountTypes")
}
func (f *fakeClient) GetAccountOutfits() ([]int, error) {
	return nil, errNotImplemented("GetAccountOutfits")
}
func (f *fakeClient) GetAccountNovelties() ([]int, error) {
	return nil, errNotImplemented("GetAccountNovelties")
}
func (f *fakeClient) GetAccountFinishers() ([]int, error) {
	return nil, errNotImplemented("GetAccountFinishers")
}
func (f *fakeClient) ResolveSkins([]int) ([]gw2api.Skin, error) {
	return nil, errNotImplemented("ResolveSkins")
}
func (f *fakeClient) ResolveColors([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveColors")
}
func (f *fakeClient) ResolveMinis([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveMinis")
}
func (f *fakeClient) ResolveMountSkins([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveMountSkins")
}
func (f *fakeClient) ResolveOutfits([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveOutfits")
}
func (f *fakeClient) ResolveNovelties([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveNovelties")
}
func (f *fakeClient) ResolveFinishers([]int) ([]gw2api.NamedEntity, error) {
	return nil, errNotImplemented("ResolveFinishers")
}
func (f *fakeClient) GetDailyAchievements() (*gw2api.DailyAchievements, error) {
	return nil, errNotImplemented("GetDailyAchievements")
}
func (f *fakeClient) GetAchievements([]int) ([]gw2api.Achievement, error) {
	return nil, errNotImplemented("GetAchievements")
}
func (f *fakeClient) GetAccountWorldBosses() ([]string, error) {
	return nil, errNotImplemented("GetAccountWorldBosses")
}
func (f *fakeClient) GetWorldBossIDs() ([]string, error) {
	return nil, errNotImplemented("GetWorldBossIDs")
}
func (f *fakeClient) GetAccountDungeons() ([]string, error) {
	return nil, errNotImplemented("GetAccountDungeons")
}
func (f *fakeClient) GetDungeons() ([]gw2api.Dungeon, error) {
	return nil, errNotImplemented("GetDungeons")
}
func (f *fakeClient) GetAccountRaids() ([]string, error) {
	return nil, errNotImplemented("GetAccountRaids")
}
func (f *fakeClient) GetRaids() ([]gw2api.Raid, error) {
	return nil, errNotImplemented("GetRaids")
}
func (f *fakeClient) GetWizardsVaultDaily() ([]gw2api.WizardsVaultObjective, error) {
	return nil, errNotImplemented("GetWizardsVaultDaily")
}
func (f *fakeClient) GetWizardsVaultWeekly() ([]gw2api.WizardsVaultObjective, error) {
	return nil, errNotImplemented("GetWizardsVaultWeekly")
}
func (f *fakeClient) GetAllAchievementIDs() ([]int, error) {
	return nil, errNotImplemented("GetAllAchievementIDs")
}
func (f *fakeClient) GetAchievementsWithProgress([]int, func(int, int, []gw2api.Achievement)) ([]gw2api.Achievement, error) {
	return nil, errNotImplemented("GetAchievementsWithProgress")
}
func (f *fakeClient) GetAccountAchievements() ([]gw2api.AccountAchievement, error) {
	return nil, errNotImplemented("GetAccountAchievements")
}
func (f *fakeClient) GetMasteryPointSummary() (*gw2api.MasteryPointSummary, error) {
	return nil, errNotImplemented("GetMasteryPointSummary")
}
func (f *fakeClient) GetLuck() (int, error) { return 0, errNotImplemented("GetLuck") }
func (f *fakeClient) GetAchievementCategories() ([]gw2api.AchievementCategory, error) {
	return nil, errNotImplemented("GetAchievementCategories")
}
func (f *fakeClient) GetAchievementGroups() ([]gw2api.AchievementGroup, error) {
	return nil, errNotImplemented("GetAchievementGroups")
}

func (f *fakeClient) resolveItems(ids []int) []gw2api.Item {
	items := make([]gw2api.Item, 0, len(ids))
	for _, id := range ids {
		if item, ok := f.items[id]; ok {
			items = append(items, item)
		}
	}
	return items
}

func errNotImplemented(method string) error {
	return fmt.Errorf("%s not implemented in fake client", method)
}
