package ui

import (
	"fmt"
	"strings"

	"gw2cli/internal/inventory"
)

func PrintGlobalHelp() {
	fmt.Println(`
GW2CLI - Guild Wars 2 Inventory Tool (v2.1.0)

Usage:
  ./gw2cli <command> [arguments]

Commands:
  search     Search your inventory for items
  list       List characters or item types
  wallet     Show account currencies
  tp         Trading Post (delivery, orders, history, price)
  exchange   Gem/Coin exchange rates
  cache      Manage local item database
  legendary  Manage Legendary Armory
  recipes    Search unlocked crafting recipes
  collection Manage account collections

Use "./gw2cli <command> -help" for more information on a command.`)
}

func PrintCollectionHelp() {
	fmt.Println(`
Usage: ./gw2cli collection [subcommand] [filter]

Subcommands:
  (default)   Show summary of all collections.
  skins       List unlocked skins.
  dyes        List unlocked dyes.
  minis       List unlocked miniatures.
  mounts      List unlocked mount skins and types.
  outfits     List unlocked outfits.
  novelties   List unlocked novelties.
  finishers   List unlocked finishers.

Arguments:
  filter      (Optional) Filter results by a partial name.`)
}

func PrintCollectionSummary(summary *inventory.CollectionSummary) {
	fmt.Println("\n--- Account Collections Summary ---")
	fmt.Printf("Skins:     %d unlocked\n", summary.Skins)
	fmt.Printf("Dyes:      %d unlocked\n", summary.Dyes)
	fmt.Printf("Minis:     %d unlocked\n", summary.Minis)
	fmt.Printf("Mounts:    %d unlocked\n", summary.Mounts)
	fmt.Printf("Outfits:   %d unlocked\n", summary.Outfits)
	fmt.Printf("Novelties: %d unlocked\n", summary.Novelties)
	fmt.Printf("Finishers: %d unlocked\n", summary.Finishers)
}

func PrintCollectionItems(items []inventory.CollectionItem, collectionType string) {
	if len(items) == 0 {
		fmt.Printf("no %s unlocked on this account\n", collectionType)
		return
	}

	fmt.Printf("\n--- Unlocked %s ---\n", collectionType)
	hasType := false
	for _, item := range items {
		if item.Type != "" {
			hasType = true
			break
		}
	}

	if hasType {
		fmt.Printf("%-40s %-20s\n", "Name", "Type")
		fmt.Println("------------------------------------------------------------")
		for _, item := range items {
			fmt.Printf("%-40s %-20s\n", item.Name, item.Type)
		}
	} else {
		fmt.Println("Name")
		fmt.Println("------------------------------------------------------------")
		for _, item := range items {
			fmt.Println(item.Name)
		}
	}
	fmt.Printf("\nTotal: %d\n", len(items))
}

func PrintRecipesHelp() {
	fmt.Println(`
Usage: ./gw2cli recipes [subcommand] [arguments]

Subcommands:
  (default)          List all unlocked crafting recipes.
  find <term>        Filter unlocked recipes by partial item name.
  ingredient <term>  Find unlocked recipes that use a specific ingredient.`)
}

func PrintRecipes(recipes []inventory.RecipeDetail) {
	if len(recipes) == 0 {
		fmt.Println("no recipes unlocked on this account")
		return
	}

	fmt.Println("\n--- Unlocked Crafting Recipes ---")
	fmt.Printf("%-35s %-25s %-5s\n", "Name", "Discipline", "Rating")
	fmt.Println("----------------------------------------------------------------------")
	for _, r := range recipes {
		fmt.Printf("%-35s %-25s %-5d\n", r.OutputName, r.Discipline, r.Rating)
	}
}

func PrintRecipeIngredientResults(recipes []inventory.RecipeDetail) {
	if len(recipes) == 0 {
		fmt.Println("no matching unlocked recipes found using this ingredient")
		return
	}

	fmt.Println("\n--- Unlocked Recipes Using Ingredient ---")
	fmt.Printf("%-30s %-20s %-7s %s\n", "Name", "Discipline", "Rating", "Ingredients")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	for _, r := range recipes {
		var ings []string
		for _, ing := range r.Ingredients {
			ings = append(ings, fmt.Sprintf("%s (x%d)", ing.Name, ing.Count))
		}
		fmt.Printf("%-30s %-20s %-7d %s\n", r.OutputName, r.Discipline, r.Rating, strings.Join(ings, ", "))
	}
}

func PrintLegendaryHelp() {
	fmt.Println(`
Usage: ./gw2cli legendary [term]

Description:
  Display all items currently in your Legendary Armory.

Arguments:
  term          (Optional) Filter results by a partial item name.`)
}

func PrintLegendaryArmory(items []inventory.LegendaryItem) {
	if len(items) == 0 {
		fmt.Println("No items found in your Legendary Armory.")
		return
	}

	fmt.Println("\n--- Account Legendary Armory ---")
	fmt.Printf("%-35s %-20s %-5s\n", "Name", "Type", "Count")
	fmt.Println("----------------------------------------------------------------")
	for _, item := range items {
		fmt.Printf("%-35s %-20s %-5d\n", item.Name, item.Type, item.Count)
	}
}

func PrintSearchHelp() {
	fmt.Println(`
Usage: ./gw2cli search [term] [flags]

Description:
  Search your entire account inventory for items.

Flags:
  -type <category>
        Filter by strict Item Type (e.g., 'Weapon', 'Armor').
  -character <name>
        Search by Character Name or Location.`)
}

func PrintListHelp() {
	fmt.Println(`
Usage: ./gw2cli list [types|characters]

Description:
  List high-level account information.

Arguments:
  types         List all unique item types in your inventory.
  characters    List all characters with details (Level, Profession, etc).`)
}

func PrintTPHelp() {
	fmt.Println(`
Usage: ./gw2cli tp <subcommand> [arguments]

Subcommands:
  delivery      Show pending Trading Post deliveries.
  orders        Show active buy and sell orders.
  history       Show past transaction history.
  price <item>  Look up current TP prices for an item (Name or ID).`)
}

func PrintExchangeHelp() {
	fmt.Println(`
Usage: ./gw2cli exchange [subcommand] [amount]

Subcommands:
  (default)     Show overview of current rates (100 gems / 100g).
  gems <amt>    How many coins can you get for <amt> gems.
  coins <amt>   How many gems can you buy with <amt> coins (in copper).`)
}

func PrintCacheHelp() {
	fmt.Println(`
Usage: ./gw2cli cache <subcommand> [arguments]

Subcommands:
  update        Build or update the local item database.
  find <term>   Search for an item in the local database by name.`)
}

func PrintResults(items []inventory.ItemDetail) {
	if len(items) == 0 {
		fmt.Println("No items found matching your criteria.")
		return
	}

	fmt.Printf("\nFound %d matching items:\n", len(items))
	for _, item := range items {
		fmt.Printf("\nName: \033[1m%s\033[0m (ID: %d)\n", item.Name, item.ID)
		fmt.Printf("Type: %s\n", item.Type)
		fmt.Printf("Total: %d\n", item.TotalCount())
		fmt.Println("Found in:")
		for _, loc := range item.Locations {
			fmt.Printf(" - %s: %s (x%d)\n", loc.Source, loc.Detail, loc.Count)
		}
	}
}

func PrintTypes(types []string) {
	fmt.Println("\n--- Available Item Types (in your inventory) ---")
	for _, t := range types {
		fmt.Println(t)
	}
}

func PrintCharacters(chars []inventory.CharacterSummary) {
	if len(chars) == 0 {
		fmt.Println("No characters found on this account.")
		return
	}

	fmt.Printf("\n--- Your Characters (%d) ---\n", len(chars))
	fmt.Printf("% -20s %-15s %-10s %-5s %-10s\n", "Name", "Profession", "Race", "Lvl", "Played")
	fmt.Println("------------------------------------------------------------------")
	for _, c := range chars {
		hours := int(c.Age.Hours())
		fmt.Printf("% -20s %-15s %-10s %-5d %dh\n", c.Name, c.Profession, c.Race, c.Level, hours)
	}
}

func PrintWallet(wallet []inventory.WalletEntry) {
	if len(wallet) == 0 {
		fmt.Println("Your wallet is empty.")
		return
	}

	fmt.Println("\n--- Account Wallet ---")
	for _, w := range wallet {
		// Format Coin (ID 1) as Gold/Silver/Copper
		if w.ID == 1 {
			fmt.Printf("%-25s %s\n", w.Name+":", FormatCoin(w.Value))
		} else {
			fmt.Printf("%-25s %d\n", w.Name+":", w.Value)
		}
	}
}

func FormatCoin(value int) string {
	gold := value / 10000
	silver := (value % 10000) / 100
	copper := value % 100
	return fmt.Sprintf("%dg %ds %dc", gold, silver, copper)
}

func PrintTPDelivery(delivery *inventory.CommerceDelivery) {
	fmt.Println("\n--- Trading Post: Ready for Pickup ---")
	fmt.Printf("Coins: %s\n", FormatCoin(delivery.Coins))
	if len(delivery.Items) > 0 {
		fmt.Println("Items:")
		for _, item := range delivery.Items {
			fmt.Printf(" - %-30s x%d\n", item.Name, item.Count)
		}
	} else {
		fmt.Println("No items waiting.")
	}
}

func PrintTPPrice(prices []inventory.CommercePrice) {
	if len(prices) == 0 {
		fmt.Println("No Trading Post data found for item.")
		return
	}

	for _, p := range prices {
		fmt.Printf("\nItem: \033[1m%s\033[0m (ID: %d)\n", p.Name, p.ID)
		fmt.Printf("Highest Buy Order: %s\n", FormatCoin(p.BuyPrice))
		fmt.Printf("Lowest Sell Offer: %s\n", FormatCoin(p.SellPrice))
	}
}

func PrintTPTransactions(buys, sells []inventory.CommerceTransaction, current bool) {
	title := "Active Orders"
	if !current {
		title = "Recent Transactions"
	}

	fmt.Printf("\n--- Trading Post: %s ---\n", title)
	
	fmt.Println("\nBUYS:")
	if len(buys) == 0 {
		fmt.Println(" None")
	} else {
		for _, tx := range buys {
			fmt.Printf(" - %-30s x%-3d at %s\n", tx.Name, tx.Quantity, FormatCoin(tx.Price))
		}
	}

	fmt.Println("\nSELLS:")
	if len(sells) == 0 {
		fmt.Println(" None")
	} else {
		for _, tx := range sells {
			fmt.Printf(" - %-30s x%-3d at %s\n", tx.Name, tx.Quantity, FormatCoin(tx.Price))
		}
	}
}

func PrintExchangeRate(gemsToCoins, coinsToGems *inventory.ExchangeRate) {
	fmt.Println("\n--- Current Gem Exchange Rates ---")
	fmt.Printf("%d Gems buys: %s (%s/gem)\n", gemsToCoins.Quantity, FormatCoin(gemsToCoins.Result), FormatCoin(gemsToCoins.CoinsPerGem))
	fmt.Printf("%s buys: %d Gems (%s/gem)\n", FormatCoin(coinsToGems.Quantity), coinsToGems.Result, FormatCoin(coinsToGems.CoinsPerGem))
}

func PrintExchangeRateSingle(rate *inventory.ExchangeRate, fromGems bool) {
	fmt.Println("\n--- Gem Exchange Rate ---")
	if fromGems {
		fmt.Printf("%d Gems buys: %s (%s/gem)\n", rate.Quantity, FormatCoin(rate.Result), FormatCoin(rate.CoinsPerGem))
	} else {
		fmt.Printf("%s buys: %d Gems (%s/gem)\n", FormatCoin(rate.Quantity), rate.Result, FormatCoin(rate.CoinsPerGem))
	}
}

func PrintCacheResults(items []inventory.CacheEntry) {
	if len(items) == 0 {
		return
	}

	for _, item := range items {
		fmt.Printf("ID: %-8d Name: %-35s Type: %s\n", item.ID, item.Name, item.Type)
	}
}
