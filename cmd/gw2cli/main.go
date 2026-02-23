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

func main() {
	flag.Usage = ui.PrintGlobalHelp
	itemFlag := flag.String("item", "", "Search by Item Name or ID")
	typeFlag := flag.String("type", "", "Filter by strict Item Type")
	charFlag := flag.String("character", "", "Search by Character Name or Location")
	listTypesFlag := flag.Bool("list-types", false, "List all unique item types")
	listCharsFlag := flag.Bool("list-characters", false, "List all characters")
	walletFlag := flag.Bool("wallet", false, "Show account wallet")
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
		case "-type", "-character", "-help":
			log.Fatalf("Error: Flag provided as value for -item. Did you forget the search term?\nUsage: ./gw2cli -item <term> [other flags]")
		}
	}

	if strings.HasPrefix(*typeFlag, "-") {
		log.Fatalf("Error: Flag provided as value for -type. Usage: ./gw2cli -type <category>")
	}
	if strings.HasPrefix(*charFlag, "-") {
		log.Fatalf("Error: Flag provided as value for -character. Usage: ./gw2cli -character <name>")
	}

	if strings.ToLower(*typeFlag) == "help" || strings.ToLower(*itemFlag) == "help" || strings.ToLower(*charFlag) == "help" {
		ui.PrintGlobalHelp()
		return
	}

	apiKey := os.Getenv("GW2_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the GW2_API_KEY environment variable")
	}
	client := gw2api.NewClient(apiKey)
	invService := inventory.NewService(client)

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
