package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"gw2cli/internal/inventory"
	"gw2cli/internal/ui"
	"gw2cli/pkg/gw2api"
)

func main() {
	// 1. Setup Flags
	flag.Usage = ui.PrintGlobalHelp
	itemFlag := flag.String("item", "", "Search by Item Name or ID")
	typeFlag := flag.String("type", "", "Filter by strict Item Type")
	charFlag := flag.String("character", "", "Search by Character Name or Location")
	listTypesFlag := flag.Bool("list-types", false, "List all unique item types")
	helpFlag := flag.Bool("help", false, "Show help")
	flag.Parse()

	// 2. Handle Immediate Help Requests
	if *helpFlag {
		ui.PrintGlobalHelp()
		return
	}
	// (Simpler check for help args)
	if strings.ToLower(*typeFlag) == "help" || strings.ToLower(*itemFlag) == "help" || strings.ToLower(*charFlag) == "help" {
		ui.PrintGlobalHelp() // Simplified for now, or route to specific help
		return
	}

	// 3. Initialize API Client
	apiKey := os.Getenv("GW2_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the GW2_API_KEY environment variable")
	}
	client := gw2api.NewClient(apiKey)

	// 4. Fetch Data (Business Logic)
	fmtStr := "Fetching account data...\n"
	os.Stdout.WriteString(fmtStr)
	
	invService := inventory.NewService(client)
	allItems, err := invService.FetchAll()
	if err != nil {
		log.Fatalf("Error fetching inventory: %v", err)
	}

	if len(allItems) == 0 {
		log.Println("Inventory is empty.")
		return
	}

	// 5. Handle -list-types
	if *listTypesFlag {
		types := inventory.GetUniqueTypes(allItems)
		ui.PrintTypes(types)
		return
	}

	// 6. Filter Results
	// Combine positional args into the "item" search if no flag provided
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

	// 7. Display Results
	ui.PrintResults(results)
}
