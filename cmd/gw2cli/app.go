package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"gw2cli/internal/inventory"
	"gw2cli/internal/ui"
	"gw2cli/pkg/gw2api"
)

const unauthenticatedCommands = "exchange, tp price, daily bosses, daily fractals, cache update, cache find, achievements update-cache"

type app struct {
	apiKey  string
	service *inventory.Service
}

type command struct {
	usage  func()
	handle func(*app, []string) error
}

func run(args []string, apiKey string) error {
	if len(args) == 0 {
		ui.PrintGlobalHelp(Version)
		return nil
	}

	client := gw2api.NewClient(apiKey)
	a := &app{
		apiKey:  apiKey,
		service: inventory.NewService(client),
	}
	a.service.Verbose = hasArg(args, "-verbose")

	commands := map[string]command{
		"account":      {usage: ui.PrintAccountHelp, handle: runAccount},
		"search":       {usage: ui.PrintSearchHelp, handle: runSearch},
		"list":         {usage: ui.PrintListHelp, handle: runList},
		"wallet":       {usage: ui.PrintWalletHelp, handle: runWallet},
		"tp":           {usage: ui.PrintTPHelp, handle: runTP},
		"legendary":    {usage: ui.PrintLegendaryHelp, handle: runLegendary},
		"recipes":      {usage: ui.PrintRecipesHelp, handle: runRecipes},
		"collection":   {usage: ui.PrintCollectionHelp, handle: runCollection},
		"daily":        {usage: ui.PrintDailyHelp, handle: runDaily},
		"weekly":       {usage: ui.PrintWeeklyHelp, handle: runWeekly},
		"achievements": {usage: ui.PrintAchievementHelp, handle: runAchievements},
		"exchange":     {usage: ui.PrintExchangeHelp, handle: runExchange},
		"cache":        {usage: ui.PrintCacheHelp, handle: runCache},
	}

	switch args[0] {
	case "-version", "--version", "version":
		fmt.Printf("GW2CLI version %s\n", Version)
		return nil
	}

	cmd, ok := commands[args[0]]
	if !ok {
		ui.PrintGlobalHelp(Version)
		return fmt.Errorf("unknown command: %s", args[0])
	}

	if len(args) > 1 && args[1] == "-help" {
		cmd.usage()
		return nil
	}

	return cmd.handle(a, args[1:])
}

func hasArg(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func (a *app) requireAPIKey() error {
	if a.apiKey != "" {
		return nil
	}
	return fmt.Errorf("GW2_API_KEY is not set; authenticated subcommands require an API key from https://account.arena.net/applications; unauthenticated subcommands: %s", unauthenticatedCommands)
}

func runAccount(a *app, args []string) error {
	accountCmd := flag.NewFlagSet("account", flag.ContinueOnError)
	if err := accountCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}

	fmt.Println("Fetching account summary...")
	acc, err := a.service.GetAccountSummary()
	if err != nil {
		return err
	}
	ui.PrintAccountSummary(acc)
	return nil
}

func runSearch(a *app, args []string) error {
	searchCmd := flag.NewFlagSet("search", flag.ContinueOnError)
	searchCmd.Usage = ui.PrintSearchHelp
	searchType := searchCmd.String("type", "", "Filter by strict Item Type")
	searchChar := searchCmd.String("character", "", "Search by Character Name or Location")
	if err := searchCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}
	if err := a.service.CheckCacheStatus(); err != nil {
		fmt.Printf("warning: %v\n", err)
	}

	fmt.Println("Fetching account data...")
	allItems, err := a.service.FetchAll()
	if err != nil {
		return err
	}
	results := inventory.Search(allItems, inventory.FilterCriteria{
		SearchTerm: strings.Join(searchCmd.Args(), " "),
		Type:       *searchType,
		Character:  *searchChar,
	})
	ui.PrintResults(results)
	return nil
}

func runList(a *app, args []string) error {
	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	listCmd.Usage = ui.PrintListHelp
	if err := listCmd.Parse(args); err != nil {
		return err
	}
	if listCmd.NArg() < 1 {
		ui.PrintListHelp()
		return nil
	}

	if err := a.requireAPIKey(); err != nil {
		return err
	}

	switch listCmd.Arg(0) {
	case "types":
		fmt.Println("Fetching account data...")
		allItems, err := a.service.FetchAll()
		if err != nil {
			return err
		}
		ui.PrintTypes(inventory.GetUniqueTypes(allItems))
	case "characters":
		fmt.Println("Fetching characters...")
		chars, err := a.service.GetCharacterList()
		if err != nil {
			return err
		}
		ui.PrintCharacters(chars)
	case "account":
		fmt.Println("Fetching account summary...")
		acc, err := a.service.GetAccountSummary()
		if err != nil {
			return err
		}
		ui.PrintAccountSummary(acc)
	default:
		ui.PrintListHelp()
	}
	return nil
}

func runWallet(a *app, args []string) error {
	walletCmd := flag.NewFlagSet("wallet", flag.ContinueOnError)
	if err := walletCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}

	fmt.Println("Fetching wallet...")
	wallet, err := a.service.GetWallet()
	if err != nil {
		return err
	}
	ui.PrintWallet(wallet)
	return nil
}

func runTP(a *app, args []string) error {
	tpCmd := flag.NewFlagSet("tp", flag.ContinueOnError)
	tpCmd.Usage = ui.PrintTPHelp
	if err := tpCmd.Parse(args); err != nil {
		return err
	}
	if tpCmd.NArg() < 1 {
		ui.PrintTPHelp()
		return nil
	}

	switch tpCmd.Arg(0) {
	case "delivery":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		delivery, err := a.service.GetDelivery()
		if err != nil {
			return err
		}
		ui.PrintTPDelivery(delivery)
	case "orders":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		buys, sells, err := a.service.GetTransactions(true)
		if err != nil {
			return err
		}
		ui.PrintTPTransactions(buys, sells, true)
	case "history":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		buys, sells, err := a.service.GetTransactions(false)
		if err != nil {
			return err
		}
		ui.PrintTPTransactions(buys, sells, false)
	case "price":
		if tpCmd.NArg() < 2 {
			return fmt.Errorf("usage: tp price <item name or ID>")
		}
		prices, err := a.service.GetPrices(strings.Join(tpCmd.Args()[1:], " "))
		if err != nil {
			return err
		}
		ui.PrintTPPrice(prices)
	default:
		ui.PrintTPHelp()
	}
	return nil
}

func runLegendary(a *app, args []string) error {
	legendaryCmd := flag.NewFlagSet("legendary", flag.ContinueOnError)
	legendaryCmd.Usage = ui.PrintLegendaryHelp
	if err := legendaryCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}
	if err := a.service.CheckCacheStatus(); err != nil {
		return err
	}

	fmt.Println("Fetching legendary armory...")
	items, err := a.service.GetLegendaryArmory(strings.Join(legendaryCmd.Args(), " "))
	if err != nil {
		return err
	}
	ui.PrintLegendaryArmory(items)
	return nil
}

func runRecipes(a *app, args []string) error {
	recipesCmd := flag.NewFlagSet("recipes", flag.ContinueOnError)
	recipesCmd.Usage = ui.PrintRecipesHelp
	if err := recipesCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}
	if err := a.service.CheckCacheStatus(); err != nil {
		return err
	}

	if recipesCmd.NArg() == 0 {
		fmt.Println("Fetching recipes...")
		recipes, err := a.service.GetUnlockedRecipes("")
		if err != nil {
			return err
		}
		ui.PrintRecipes(recipes)
		return nil
	}

	switch recipesCmd.Arg(0) {
	case "find":
		if recipesCmd.NArg() < 2 {
			return fmt.Errorf("usage: recipes find <term>")
		}
		fmt.Println("Searching recipes...")
		recipes, err := a.service.GetUnlockedRecipes(strings.Join(recipesCmd.Args()[1:], " "))
		if err != nil {
			return err
		}
		ui.PrintRecipes(recipes)
	case "ingredient":
		if recipesCmd.NArg() < 2 {
			return fmt.Errorf("usage: recipes ingredient <term>")
		}
		fmt.Println("Searching recipes by ingredient...")
		recipes, err := a.service.SearchRecipesByIngredient(strings.Join(recipesCmd.Args()[1:], " "))
		if err != nil {
			return err
		}
		ui.PrintRecipeIngredientResults(recipes)
	default:
		ui.PrintRecipesHelp()
	}
	return nil
}

func runCollection(a *app, args []string) error {
	collectionCmd := flag.NewFlagSet("collection", flag.ContinueOnError)
	collectionCmd.Usage = ui.PrintCollectionHelp
	if err := collectionCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}

	if collectionCmd.NArg() == 0 {
		fmt.Println("Fetching collection summary...")
		summary, err := a.service.GetCollectionSummary()
		if err != nil {
			return err
		}
		ui.PrintCollectionSummary(summary)
		return nil
	}

	sub := collectionCmd.Arg(0)
	filter := ""
	if collectionCmd.NArg() > 1 {
		filter = strings.Join(collectionCmd.Args()[1:], " ")
	}

	fmt.Printf("Fetching %s...\n", sub)
	var items []inventory.CollectionItem
	var err error
	switch sub {
	case "skins":
		items, err = a.service.GetCollectionSkins(filter)
	case "dyes":
		items, err = a.service.GetCollectionDyes(filter)
	case "minis":
		items, err = a.service.GetCollectionMinis(filter)
	case "mounts":
		items, err = a.service.GetCollectionMounts(filter)
	case "outfits":
		items, err = a.service.GetCollectionOutfits(filter)
	case "novelties":
		items, err = a.service.GetCollectionNovelties(filter)
	case "finishers":
		items, err = a.service.GetCollectionFinishers(filter)
	default:
		ui.PrintCollectionHelp()
		return nil
	}
	if err != nil {
		return err
	}
	ui.PrintCollectionItems(items, sub)
	return nil
}

func runDaily(a *app, args []string) error {
	dailyCmd := flag.NewFlagSet("daily", flag.ContinueOnError)
	dailyCmd.Usage = ui.PrintDailyHelp
	if err := dailyCmd.Parse(args); err != nil {
		return err
	}

	if dailyCmd.NArg() == 0 {
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		fmt.Println("Fetching daily status...")
		bosses, err := a.service.GetDailyBosses()
		if err != nil {
			return err
		}
		dungeons, err := a.service.GetDailyDungeons()
		if err != nil {
			return err
		}
		ui.PrintDailyStatus(bosses, dungeons)

		fractals, err := a.service.GetDailyFractals()
		if err != nil {
			return err
		}
		ui.PrintFractalDailies(fractals)

		wv, err := a.service.GetDailyWizardsVault()
		if err != nil {
			return err
		}
		ui.PrintWizardsVault(wv, "Daily")
		return nil
	}

	switch dailyCmd.Arg(0) {
	case "fractals":
		fractals, err := a.service.GetDailyFractals()
		if err != nil {
			return err
		}
		ui.PrintFractalDailies(fractals)
	case "bosses":
		bosses, err := a.dailyBosses()
		if err != nil {
			return err
		}
		ui.PrintDailyStatus(bosses, nil)
	case "wizardsvault":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		wv, err := a.service.GetDailyWizardsVault()
		if err != nil {
			return err
		}
		ui.PrintWizardsVault(wv, "Daily")
	default:
		ui.PrintDailyHelp()
	}
	return nil
}

func (a *app) dailyBosses() ([]inventory.DailyStatus, error) {
	if a.apiKey == "" {
		return a.service.GetWorldBossList()
	}
	return a.service.GetDailyBosses()
}

func runWeekly(a *app, args []string) error {
	weeklyCmd := flag.NewFlagSet("weekly", flag.ContinueOnError)
	weeklyCmd.Usage = ui.PrintWeeklyHelp
	if err := weeklyCmd.Parse(args); err != nil {
		return err
	}
	if err := a.requireAPIKey(); err != nil {
		return err
	}

	fmt.Println("Fetching weekly status...")
	raids, err := a.service.GetWeeklyRaids()
	if err != nil {
		return err
	}
	ui.PrintRaidStatus(raids)

	wv, err := a.service.GetWeeklyWizardsVault()
	if err != nil {
		return err
	}
	ui.PrintWizardsVault(wv, "Weekly")
	return nil
}

func runAchievements(a *app, args []string) error {
	achievementsCmd := flag.NewFlagSet("achievements", flag.ContinueOnError)
	achievementsCmd.Usage = ui.PrintAchievementHelp
	status := achievementsCmd.String("status", "any", "Filter by status (any, completed, incomplete)")
	if err := achievementsCmd.Parse(args); err != nil {
		return err
	}

	if achievementsCmd.NArg() == 0 {
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		if err := a.service.CheckAchievementCacheStatus(); err != nil {
			return err
		}
		fmt.Println("Fetching achievement summary...")
		summaries, err := a.service.GetAchievementSummary()
		if err != nil {
			return err
		}
		ui.PrintAchievementSummary(summaries)
		return nil
	}

	switch achievementsCmd.Arg(0) {
	case "update-cache":
		return a.service.EnsureAchievementCache(true)
	case "all":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		if err := a.service.CheckAchievementCacheStatus(); err != nil {
			return err
		}
		fmt.Println("Fetching all achievements...")
		results, err := a.service.GetAllAchievements(*status)
		if err != nil {
			return err
		}
		ui.PrintAchievementTable(results)
	case "find":
		if achievementsCmd.NArg() < 2 {
			return fmt.Errorf("usage: achievements find <term>")
		}
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		if err := a.service.CheckAchievementCacheStatus(); err != nil {
			return err
		}
		fmt.Println("Searching achievements...")
		results, err := a.service.FindAchievements(strings.Join(achievementsCmd.Args()[1:], " "))
		if err != nil {
			return err
		}
		ui.PrintAchievementProgress(results)
	case "masteries":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		fmt.Println("Fetching mastery summary...")
		summary, err := a.service.GetMasterySummary()
		if err != nil {
			return err
		}
		ui.PrintMasterySummary(summary)
	case "convergences", "raids", "fractals", "strikes", "pvp", "wvw":
		if err := a.requireAPIKey(); err != nil {
			return err
		}
		if err := a.service.CheckAchievementCacheStatus(); err != nil {
			return err
		}
		category := achievementsCmd.Arg(0)
		fmt.Printf("Fetching %s achievements...\n", category)
		results, err := a.service.GetCategoryAchievements(category)
		if err != nil {
			return err
		}
		ui.PrintAchievementProgress(results)
	default:
		ui.PrintAchievementHelp()
	}
	return nil
}

func runExchange(a *app, args []string) error {
	exchangeCmd := flag.NewFlagSet("exchange", flag.ContinueOnError)
	exchangeCmd.Usage = ui.PrintExchangeHelp
	if err := exchangeCmd.Parse(args); err != nil {
		return err
	}

	if exchangeCmd.NArg() == 0 {
		g2c, err := a.service.GetExchangeRate(100, true)
		if err != nil {
			return err
		}
		c2g, err := a.service.GetExchangeRate(1000000, false)
		if err != nil {
			return err
		}
		ui.PrintExchangeRate(g2c, c2g)
		return nil
	}

	if exchangeCmd.NArg() < 2 {
		ui.PrintExchangeHelp()
		return nil
	}

	amount, err := strconv.Atoi(exchangeCmd.Arg(1))
	if err != nil || amount <= 0 {
		return fmt.Errorf("amount must be a positive integer")
	}

	switch exchangeCmd.Arg(0) {
	case "gems":
		rate, err := a.service.GetExchangeRate(amount, true)
		if err != nil {
			return err
		}
		ui.PrintExchangeRateSingle(rate, true)
	case "coins":
		rate, err := a.service.GetExchangeRate(amount, false)
		if err != nil {
			return err
		}
		ui.PrintExchangeRateSingle(rate, false)
	default:
		ui.PrintExchangeHelp()
	}
	return nil
}

func runCache(a *app, args []string) error {
	cacheCmd := flag.NewFlagSet("cache", flag.ContinueOnError)
	cacheCmd.Usage = ui.PrintCacheHelp
	if err := cacheCmd.Parse(args); err != nil {
		return err
	}
	if cacheCmd.NArg() < 1 {
		ui.PrintCacheHelp()
		return nil
	}

	switch cacheCmd.Arg(0) {
	case "update":
		return a.service.EnsureCache(true)
	case "find":
		if cacheCmd.NArg() < 2 {
			return fmt.Errorf("usage: cache find <term>")
		}
		matches, err := a.service.FindInCache(strings.Join(cacheCmd.Args()[1:], " "))
		if err != nil {
			return err
		}
		ui.PrintCacheResults(matches)
	default:
		ui.PrintCacheHelp()
	}
	return nil
}
