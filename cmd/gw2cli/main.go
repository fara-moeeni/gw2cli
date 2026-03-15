package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gw2cli/internal/inventory"
	"gw2cli/internal/ui"
	"gw2cli/pkg/gw2api"
)

const Version = "1.5.0"

func main() {
        flag.Usage = ui.PrintGlobalHelp
        itemFlag := flag.String("item", "", "Search by Item Name or ID")
        typeFlag := flag.String("type", "", "Filter by strict Item Type")
        charFlag := flag.String("character", "", "Search by Character Name or Location")
        listTypesFlag := flag.Bool("list-types", false, "List all unique item types")
        listCharsFlag := flag.Bool("list-characters", false, "List all characters")
        walletFlag := flag.Bool("wallet", false, "Show account wallet")
        tpDeliveryFlag := flag.Bool("tp-delivery", false, "Show pending Trading Post deliveries")
        tpOrdersFlag := flag.Bool("tp-orders", false, "Show active Trading Post orders")
        tpHistoryFlag := flag.Bool("tp-history", false, "Show past Trading Post transactions")
        tpPriceFlag := flag.String("tp-price", "", "Check current Trading Post prices")
        exchangeFlag := flag.Bool("exchange", false, "Show current gem/coin exchange rates")
        exchangeCoinsFlag := flag.Int("exchange-coins", 0, "Amount of coins to exchange for gems")
        exchangeGemsFlag := flag.Int("exchange-gems", 0, "Amount of gems to exchange for coins")
        updateCacheFlag := flag.Bool("update-cache", false, "Update local item database")
        findFlag := flag.String("find", "", "Search for items in local database by name")

        verboseFlag := flag.Bool("verbose", false, "Enable verbose output")
        helpFlag := flag.Bool("help", false, "Show help")
        flag.Parse()

        if *helpFlag {
                ui.PrintGlobalHelp()
                return
        }

        // Correct flag parsing edge cases where flags are consumed as values
        if strings.HasPrefix(*itemFlag, "-") {
                switch *itemFlag {
                case "-list-types":
                        *listTypesFlag = true
                        *itemFlag = ""
                case "-list-characters":
                        *listCharsFlag = true
                        *itemFlag = ""
                case "-wallet":
                        *walletFlag = true
                        *itemFlag = ""
                case "-tp-delivery":
                        *tpDeliveryFlag = true
                        *itemFlag = ""
                case "-tp-orders":
                        *tpOrdersFlag = true
                        *itemFlag = ""
                case "-tp-history":
                        *tpHistoryFlag = true
                        *itemFlag = ""
                case "-type", "-character", "-tp-price", "-verbose", "-help", "-update-cache", "-find":
                        log.Fatalf("Error: Flag provided as value for -item. Did you forget the search term?\nUsage: ./gw2cli -item <term> [other flags]")
                }
        }

        if strings.HasPrefix(*typeFlag, "-") {
                log.Fatalf("Error: Flag provided as value for -type. Usage: ./gw2cli -type <category>")
        }
        if strings.HasPrefix(*charFlag, "-") {
                log.Fatalf("Error: Flag provided as value for -character. Usage: ./gw2cli -character <name>")
        }
        if strings.HasPrefix(*tpPriceFlag, "-") {
                log.Fatalf("Error: Flag provided as value for -tp-price. Usage: ./gw2cli -tp-price <name or ID>")
        }
        if strings.HasPrefix(*findFlag, "-") {
                log.Fatalf("Error: Flag provided as value for -find. Usage: ./gw2cli -find <term>")
        }

        if strings.ToLower(*typeFlag) == "help" || strings.ToLower(*itemFlag) == "help" || strings.ToLower(*charFlag) == "help" {
                ui.PrintGlobalHelp()
                return
        }

        apiKey := os.Getenv("GW2_API_KEY")
        client := gw2api.NewClient(apiKey)
        invService := inventory.NewService(client)
        invService.Verbose = *verboseFlag

        // Global cache status check
        _ = invService.CheckCacheStatus()

        if *updateCacheFlag {
                if err := invService.EnsureCache(true); err != nil {
                        log.Fatalf("Error updating cache: %v", err)
                }
                return
        }

        if *findFlag != "" {
                matches, err := invService.FindInCache(*findFlag)
                if err != nil {
                        log.Fatalf("Error searching cache: %v", err)
                }
                if len(matches) == 0 {
                        fmt.Printf("no items found matching \"%s\"\n", *findFlag)
                        return
                }
                ui.PrintCacheResults(matches)
                return
        }

        // Auth-required check
        if apiKey == "" {
                log.Fatal("Please set the GW2_API_KEY environment variable")
        }
	if *listCharsFlag {
		fmt.Println("Fetching character list...")
		chars, err := invService.GetCharacterList()
		if err != nil {
			log.Fatalf("Error fetching characters: %v", err)
		}
		ui.PrintCharacters(chars)
		return
	}

	        if *walletFlag {
	                fmt.Println("Fetching account wallet...")
	                wallet, err := invService.GetWallet()
	                if err != nil {
	                        log.Fatalf("Error fetching wallet: %v", err)
	                }
	                ui.PrintWallet(wallet)
	                return
	        }
	
	        if *tpDeliveryFlag {
	                fmt.Println("Fetching TP delivery...")
	                delivery, err := invService.GetDelivery()
	                if err != nil {
	                        log.Fatalf("Error fetching TP delivery: %v", err)
	                }
	                ui.PrintTPDelivery(delivery)
	                return
	        }
	
	        if *tpOrdersFlag {
	                fmt.Println("Fetching active TP orders...")
	                buys, sells, err := invService.GetTransactions(true)
	                if err != nil {
	                        log.Fatalf("Error fetching TP orders: %v", err)
	                }
	                ui.PrintTPTransactions(buys, sells, true)
	                return
	        }
	
	        if *tpHistoryFlag {
	                fmt.Println("Fetching TP transaction history...")
	                buys, sells, err := invService.GetTransactions(false)
	                if err != nil {
	                        log.Fatalf("Error fetching TP history: %v", err)
	                }
	                ui.PrintTPTransactions(buys, sells, false)
	                return
	        }
	
	                        if *tpPriceFlag != "" {
	                                fmt.Println("Fetching TP prices...")
	                                prices, err := invService.GetPrices(*tpPriceFlag)
	                                if err != nil {
	                                        log.Fatalf("Error fetching TP prices: %v", err)
	                                }
	                                ui.PrintTPPrice(prices)
	                                return
	                        }
	        
	                        if *exchangeFlag {
	                                fmt.Println("Fetching current exchange rates...")
	                                g2c, err := invService.GetExchangeRate(100, true)
	                                if err != nil {
	                                        log.Fatalf("Error fetching gem exchange: %v", err)
	                                }
	                                c2g, err := invService.GetExchangeRate(1000000, false) // 100g
	                                if err != nil {
	                                        log.Fatalf("Error fetching coin exchange: %v", err)
	                                }
	                                ui.PrintExchangeRate(g2c, c2g)
	                                return
	                        }
	        
	                        if *exchangeGemsFlag > 0 {
	                                fmt.Printf("Calculating exchange for %d gems...\n", *exchangeGemsFlag)
	                                rate, err := invService.GetExchangeRate(*exchangeGemsFlag, true)
	                                if err != nil {
	                                        log.Fatalf("Error fetching exchange: %v", err)
	                                }
	                                ui.PrintExchangeRateSingle(rate, true)
	                                return
	                        }
	        
	                        if *exchangeCoinsFlag > 0 {
	                                fmt.Printf("Calculating exchange for %s...\n", ui.FormatCoin(*exchangeCoinsFlag))
	                                rate, err := invService.GetExchangeRate(*exchangeCoinsFlag, false)
	                                if err != nil {
	                                        log.Fatalf("Error fetching exchange: %v", err)
	                                }
	                                ui.PrintExchangeRateSingle(rate, false)
	                                return
	                        }
	        
	                        os.Stdout.WriteString("Fetching account data...\n")
	        	allItems, err := invService.FetchAll()
	if err != nil {
		log.Fatalf("Error fetching inventory: %v", err)
	}

	if len(allItems) == 0 {
		log.Println("Inventory is empty.")
		return
	}

	if *listTypesFlag {
		types := inventory.GetUniqueTypes(allItems)
		ui.PrintTypes(types)
		return
	}

	// Combine positional args into the search term if needed
	globalFilter := strings.Join(flag.Args(), " ")
	searchTerm := *itemFlag
	if searchTerm == "" {
		searchTerm = globalFilter
	} else if globalFilter != "" {
		searchTerm += " " + globalFilter
	}

	criteria := inventory.FilterCriteria{
		SearchTerm: searchTerm,
		Type:       *typeFlag,
		Character:  *charFlag,
	}

	results := inventory.Search(allItems, criteria)
	ui.PrintResults(results)
}
