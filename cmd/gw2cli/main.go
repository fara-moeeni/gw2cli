package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gw2cli/internal/inventory"
	"gw2cli/internal/ui"
	"gw2cli/pkg/gw2api"
)

const Version = "2.2.0"

func main() {
	if len(os.Args) < 2 {
		ui.PrintGlobalHelp()
		return
	}

	// Define Subcommands
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchType := searchCmd.String("type", "", "Filter by strict Item Type")
	searchChar := searchCmd.String("character", "", "Search by Character Name or Location")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
	tpCmd := flag.NewFlagSet("tp", flag.ExitOnError)
	exchangeCmd := flag.NewFlagSet("exchange", flag.ExitOnError)
	cacheCmd := flag.NewFlagSet("cache", flag.ExitOnError)
	legendaryCmd := flag.NewFlagSet("legendary", flag.ExitOnError)
	recipesCmd := flag.NewFlagSet("recipes", flag.ExitOnError)

	// Configure usages
	searchCmd.Usage = ui.PrintSearchHelp
	listCmd.Usage = ui.PrintListHelp
	tpCmd.Usage = ui.PrintTPHelp
	exchangeCmd.Usage = ui.PrintExchangeHelp
	cacheCmd.Usage = ui.PrintCacheHelp
	legendaryCmd.Usage = ui.PrintLegendaryHelp
	recipesCmd.Usage = ui.PrintRecipesHelp

	verbose := false
	for _, arg := range os.Args {
		if arg == "-verbose" {
			verbose = true
			break
		}
	}

	apiKey := os.Getenv("GW2_API_KEY")
	client := gw2api.NewClient(apiKey)
	invService := inventory.NewService(client)
	invService.Verbose = verbose

	switch os.Args[1] {
	case "search":
		searchCmd.Parse(os.Args[2:])
		requireAPIKey(apiKey)
		_ = invService.CheckCacheStatus()

		searchTerm := strings.Join(searchCmd.Args(), " ")
		fmt.Println("Fetching account data...")
		allItems, err := invService.FetchAll()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		results := inventory.Search(allItems, inventory.FilterCriteria{
			SearchTerm: searchTerm,
			Type:       *searchType,
			Character:  *searchChar,
		})
		ui.PrintResults(results)

	case "list":
		listCmd.Parse(os.Args[2:])
		if listCmd.NArg() < 1 {
			ui.PrintListHelp()
			return
		}
		switch listCmd.Arg(0) {
		case "types":
			requireAPIKey(apiKey)
			fmt.Println("Fetching account data...")
			allItems, err := invService.FetchAll()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintTypes(inventory.GetUniqueTypes(allItems))
		case "characters":
			requireAPIKey(apiKey)
			fmt.Println("Fetching characters...")
			chars, err := invService.GetCharacterList()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintCharacters(chars)
		default:
			ui.PrintListHelp()
		}

	case "wallet":
		walletCmd.Parse(os.Args[2:])
		requireAPIKey(apiKey)
		fmt.Println("Fetching wallet...")
		wallet, err := invService.GetWallet()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		ui.PrintWallet(wallet)

	case "tp":
		tpCmd.Parse(os.Args[2:])
		if tpCmd.NArg() < 1 {
			ui.PrintTPHelp()
			return
		}
		switch tpCmd.Arg(0) {
		case "delivery":
			requireAPIKey(apiKey)
			delivery, err := invService.GetDelivery()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintTPDelivery(delivery)
		case "orders":
			requireAPIKey(apiKey)
			buys, sells, err := invService.GetTransactions(true)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintTPTransactions(buys, sells, true)
		case "history":
			requireAPIKey(apiKey)
			buys, sells, err := invService.GetTransactions(false)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintTPTransactions(buys, sells, false)
		case "price":
			if tpCmd.NArg() < 2 {
				log.Fatal("Usage: tp price <item name or ID>")
			}
			prices, err := invService.GetPrices(strings.Join(tpCmd.Args()[1:], " "))
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintTPPrice(prices)
		default:
			ui.PrintTPHelp()
		}

	case "legendary":
		legendaryCmd.Parse(os.Args[2:])
		requireAPIKey(apiKey)
		if err := invService.CheckCacheStatus(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		term := strings.Join(legendaryCmd.Args(), " ")
		fmt.Println("Fetching legendary armory...")
		items, err := invService.GetLegendaryArmory(term)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		ui.PrintLegendaryArmory(items)

	case "recipes":
		recipesCmd.Parse(os.Args[2:])
		requireAPIKey(apiKey)
		if err := invService.CheckCacheStatus(); err != nil {
			log.Fatalf("Error: %v", err)
		}

		if recipesCmd.NArg() == 0 {
			fmt.Println("Fetching recipes...")
			recipes, err := invService.GetUnlockedRecipes("")
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintRecipes(recipes)
		} else {
			switch recipesCmd.Arg(0) {
			case "find":
				if recipesCmd.NArg() < 2 {
					log.Fatal("Usage: recipes find <term>")
				}
				fmt.Println("Searching recipes...")
				recipes, err := invService.GetUnlockedRecipes(strings.Join(recipesCmd.Args()[1:], " "))
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				ui.PrintRecipes(recipes)
			case "ingredient":
				if recipesCmd.NArg() < 2 {
					log.Fatal("Usage: recipes ingredient <term>")
				}
				fmt.Println("Searching recipes by ingredient...")
				recipes, err := invService.SearchRecipesByIngredient(strings.Join(recipesCmd.Args()[1:], " "))
				if err != nil {
					log.Fatalf("Error: %v", err)
				}
				ui.PrintRecipeIngredientResults(recipes)
			default:
				ui.PrintRecipesHelp()
			}
		}

	case "exchange":
		exchangeCmd.Parse(os.Args[2:])
		if exchangeCmd.NArg() == 0 {
			g2c, _ := invService.GetExchangeRate(100, true)
			c2g, _ := invService.GetExchangeRate(1000000, false)
			ui.PrintExchangeRate(g2c, c2g)
		} else {
			if exchangeCmd.NArg() < 2 {
				ui.PrintExchangeHelp()
				return
			}
			amount, _ := strconv.Atoi(exchangeCmd.Arg(1))
			switch exchangeCmd.Arg(0) {
			case "gems":
				rate, _ := invService.GetExchangeRate(amount, true)
				ui.PrintExchangeRateSingle(rate, true)
			case "coins":
				rate, _ := invService.GetExchangeRate(amount, false)
				ui.PrintExchangeRateSingle(rate, false)
			default:
				ui.PrintExchangeHelp()
			}
		}

	case "cache":
		cacheCmd.Parse(os.Args[2:])
		if cacheCmd.NArg() < 1 {
			ui.PrintCacheHelp()
			return
		}
		switch cacheCmd.Arg(0) {
		case "update":
			if err := invService.EnsureCache(true); err != nil {
				log.Fatalf("Error: %v", err)
			}
		case "find":
			if cacheCmd.NArg() < 2 {
				log.Fatal("Usage: cache find <term>")
			}
			matches, err := invService.FindInCache(strings.Join(cacheCmd.Args()[1:], " "))
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			ui.PrintCacheResults(matches)
		default:
			ui.PrintCacheHelp()
		}

	case "-version", "--version", "version":
		fmt.Printf("GW2CLI version %s\n", Version)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		ui.PrintGlobalHelp()
		os.Exit(1)
	}
}

func requireAPIKey(key string) {
	if key == "" {
		log.Fatal("Error: GW2_API_KEY environment variable not set")
	}
}
