package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gw2cli/internal/inventory"

	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintGlobalHelp() {
	fmt.Println(`
GW2CLI - Guild Wars 2 Inventory Tool (v2.8.1)

Usage:
  ./gw2cli <command> [arguments]

Commands:
  account    Show account summary (Fractal Level, AP, etc)
  search     Search your inventory for items
  list       List characters or item types
  wallet     Show account currencies
  tp         Trading Post (delivery, orders, history, price)
  exchange   Gem/Coin exchange rates
  cache      Manage local item database
  legendary  Manage Legendary Armory
  recipes    Search unlocked crafting recipes
  collection Manage account collections
  daily      Track daily resets and objectives
  weekly     Track weekly resets and raids
  achievements Track achievements and masteries

Use "./gw2cli <command> -help" for more information on a command.`)
}

func PrintAccountSummary(acc *inventory.AccountSummary) {
	fmt.Println("\n--- Account Summary ---")
	fmt.Printf("Account Name:  %s\n", acc.Name)
	fmt.Printf("Fractal Level: %d\n", acc.FractalLevel)
	fmt.Printf("WvW Rank:      %d\n", acc.WvwRank)
	fmt.Printf("Daily AP:      %d\n", acc.DailyAP)
	fmt.Printf("Monthly AP:    %d\n", acc.MonthlyAP)
	
	created, _ := time.Parse(time.RFC3339, acc.Created)
	fmt.Printf("Created:       %s\n", created.Format("2006-01-02"))
	
	// Calculate actual Account Age (time since creation)
	elapsed := time.Since(created)
	accYears := int(elapsed.Hours() / 24 / 365)
	accDays := int(elapsed.Hours()/24) % 365
	fmt.Printf("Account Age:   %d years, %d days\n", accYears, accDays)

	// Playtime calculation
	playYears := acc.TotalPlaytime / (365 * 24 * 3600)
	playDays := (acc.TotalPlaytime % (365 * 24 * 3600)) / (24 * 3600)
	playHours := (acc.TotalPlaytime % (24 * 3600)) / 3600
	fmt.Printf("Total Playtime: %d years, %d days, %d hours\n", playYears, playDays, playHours)
}

func PrintAchievementHelp() {
	fmt.Println(`
Usage: ./gw2cli achievements [subcommand] [filter]

Subcommands:
  (default)           Summary of achievement categories and AP.
  all [--status=...]  List all achievements (filters: completed, incomplete, any).
  update-cache        Build/refresh the local achievement database.
  find <term>         Search achievements by name and show progress.
  masteries           Show mastery points and luck progress.
  convergences        Convergences achievement progress.
  raids               Raid achievement progress.
  fractals            Fractal achievement progress.
  strikes             Strike mission achievement progress.
  pvp                 PvP rank and achievement progress.
  wvw                 WvW rank and achievement progress.`)
}

func PrintAchievementTable(achievements []inventory.AchievementProgress) {
	if len(achievements) == 0 {
		fmt.Println("No matching achievements found.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Name", "Progress", "Status", "AP"})

	for _, a := range achievements {
		prog := fmt.Sprintf("%d/%d", a.Current, a.Max)
		if a.Max == 0 {
			prog = fmt.Sprintf("%d", a.Current)
		}
		t.AppendRow(table.Row{
			a.ID,
			a.Name,
			prog,
			a.StatusSymbol,
			a.Points,
		})
	}
	t.SetStyle(table.StyleLight)
	t.Render()
}

func PrintAchievementSummary(summaries []inventory.CategorySummary) {
	fmt.Println("\n--- Achievement Summary ---")
	fmt.Printf("%-40s %-15s %-10s\n", "Category", "Completed", "AP Earned")
	fmt.Println("----------------------------------------------------------------------")
	for _, s := range summaries {
		fmt.Printf("%-40s %d/%-13d %-10d\n", s.Name, s.Completed, s.Total, s.AP)
	}
}

func PrintAchievementProgress(achievements []inventory.AchievementProgress) {
	if len(achievements) == 0 {
		fmt.Println("No matching achievements found.")
		return
	}

	fmt.Println("\n--- Achievement Progress ---")
	fmt.Printf("%-40s %-12s %-12s %-5s %s\n", "Name", "Progress", "Tiers", "AP", "Status")
	fmt.Println("-----------------------------------------------------------------------------------")
	for _, a := range achievements {
		prog := fmt.Sprintf("%d/%d", a.Current, a.Max)
		if a.Max == 0 {
			prog = fmt.Sprintf("%d", a.Current)
		}
		fmt.Printf("%-40s %-12s %-12s %-5d %s\n", a.Name, prog, a.TierStatus, a.Points, a.StatusSymbol)
	}
}

func PrintMasterySummary(summary *inventory.MasterySummary) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Region", "Spent", "Earned"})

	total := 0
	var summaryParts []string
	for _, r := range summary.Regions {
		t.AppendRow(table.Row{r.Name, r.Spent, r.Earned})
		total += r.Earned
		summaryParts = append(summaryParts, fmt.Sprintf("%s: %d", r.Name, r.Earned))
	}
	t.SetStyle(table.StyleLight)
	t.Render()

	fmt.Printf("\n%s  Total: %d\n", strings.Join(summaryParts, "  "), total)
	fmt.Printf("Total Luck: %d\n", summary.Luck)
}

func PrintDailyHelp() {
	fmt.Println(`
Usage: ./gw2cli daily [subcommand]

Description:
  Track daily reset status for various account activities.

Subcommands:
  (default)     Summary of bosses, dungeons, fractals, and Wizard's Vault.
  fractals      Today's fractal dailies and recommended.
  bosses        World boss daily reset status.
  wizardsvault  Daily Wizard's Vault objectives and progress.`)
}

func PrintWeeklyHelp() {
	fmt.Println(`
Usage: ./gw2cli weekly

Description:
  Track weekly reset status for raids and Wizard's Vault objectives.`)
}

func PrintDailyStatus(bosses, dungeons []inventory.DailyStatus) {
	fmt.Println("\n--- Daily Reset Status ---")
	fmt.Println("Daily reset: 00:00 UTC")

	fmt.Println("\nWORLD BOSSES:")
	for _, b := range bosses {
		status := "[ ]"
		if b.Completed {
			status = "[✓]"
		}
		fmt.Printf(" %s %s\n", status, b.Name)
	}

	fmt.Println("\nDUNGEONS:")
	for _, d := range dungeons {
		status := "[ ]"
		if d.Completed {
			status = "[✓]"
		}
		fmt.Printf(" %s %s\n", status, d.Name)
	}
}

func PrintFractalDailies(fractals []inventory.FractalDaily) {
	fmt.Println("\n--- Daily Fractals ---")
	fmt.Println("Daily reset: 00:00 UTC")

	for _, f := range fractals {
		status := "[ ]"
		if f.Completed {
			status = "[✓]"
		}
		fmt.Printf(" %s %-5s %s\n", status, f.Tier, f.Name)
	}
}

func PrintWizardsVault(objectives []inventory.WizardsVaultStatus, title string) {
	fmt.Printf("\n--- Wizard's Vault: %s ---\n", title)
	if title == "Weekly" {
		fmt.Println("Weekly reset: Tuesday 00:00 UTC")
	} else {
		fmt.Println("Daily reset: 00:00 UTC")
	}

	fmt.Printf("\n%-45s %-10s %-8s %s\n", "Objective", "Progress", "Acclaim", "Status")
	fmt.Println("----------------------------------------------------------------------------")
	for _, o := range objectives {
		status := "[ ]"
		if o.Completed {
			status = "[✓]"
		}
		prog := fmt.Sprintf("%d/%d", o.ProgressCur, o.ProgressGoal)
		fmt.Printf("%-45s %-10s %-8d %s\n", o.Title, prog, o.Acclaim, status)
	}
}

func PrintRaidStatus(wings []inventory.RaidWingStatus) {
	fmt.Println("\n--- Weekly Raid Status ---")
	fmt.Println("Weekly reset: Tuesday 00:00 UTC")

	for _, wing := range wings {
		fmt.Printf("\n%s:\n", wing.Name)
		for _, event := range wing.Events {
			status := "[ ]"
			if event.Completed {
				status = "[✓]"
			}
			fmt.Printf(" %s %s\n", status, event.Name)
		}
	}
	
	fmt.Println("\nOTHERS:")
	fmt.Println(" Convergences: rotation not available via API — check https://wiki.guildwars2.com/wiki/Convergences")
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
Usage: ./gw2cli list [types|characters|account]

Description:
  List high-level account information.

Arguments:
  types         List all unique item types in your inventory.
  characters    List all characters with details (Level, Profession, etc).
  account       Show general account summary (Fractal Level, AP, etc).`)
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
