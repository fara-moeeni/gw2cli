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
	// 1. Setup Flags
	flag.Usage = ui.PrintGlobalHelp
	itemFlag := flag.String("item", "", "Search by Item Name or ID")
	typeFlag := flag.String("type", "", "Filter by strict Item Type")
	charFlag := flag.String("character", "", "Search by Character Name or Location")
	listTypesFlag := flag.Bool("list-types", false, "List all unique item types")
	listCharsFlag := flag.Bool("list-characters", false, "List all characters")
	helpFlag := flag.Bool("help", false, "Show help")
	flag.Parse()

	// 2. Handle Immediate Help Requests
	if *helpFlag {
		ui.PrintGlobalHelp()
		return
	}

	// Safety Check: Did the user type "-item -list-types"?
	// If *itemFlag starts with "-", the parser consumed the next flag as the value.
	if strings.HasPrefix(*itemFlag, "-") {
		// If they explicitly wanted to search for a string starting with "-", they can escape it.
		// But 99% of the time, this is a typo.
		// Check if the "value" matches a known flag name
		if *itemFlag == "-list-types" {
			*listTypesFlag = true
			*itemFlag = "" // Reset item flag so we don't search for it
		} else if *itemFlag == "-list-characters" {
			*listCharsFlag = true
			*itemFlag = ""
		} else if *itemFlag == "-type" || *itemFlag == "-character" || *itemFlag == "-help" {
			log.Fatalf("Error: Flag provided as value for -item. Did you forget the search term?\nUsage: ./gw2cli -item <term> [other flags]")
		}
	}

	if strings.HasPrefix(*typeFlag, "-") {
		log.Fatalf("Error: Flag provided as value for -type. Usage: ./gw2cli -type <category>")
	}
	if strings.HasPrefix(*charFlag, "-") {
		log.Fatalf("Error: Flag provided as value for -character. Usage: ./gw2cli -character <name>")
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
	invService := inventory.NewService(client)

	// 4. Handle -list-characters (Exclusive Action)
	if *listCharsFlag {
		fmt.Println("Fetching character list...")
		chars, err := invService.GetCharacterList()
		if err != nil {
			log.Fatalf("Error fetching characters: %v", err)
		}
		ui.PrintCharacters(chars)
		return
	}

	// 5. Fetch Data (Inventory)
	fmtStr := "Fetching account data...\n"
	os.Stdout.WriteString(fmtStr)
	
	allItems, err := invService.FetchAll()
	if err != nil {
		log.Fatalf("Error fetching inventory: %v", err)
	}

	if len(allItems) == 0 {
		log.Println("Inventory is empty.")
		return
	}

	// 6. Handle -list-types
	if *listTypesFlag {
		types := inventory.GetUniqueTypes(allItems)
		ui.PrintTypes(types)
		return
	}

	// 7. Filter Results
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

	// 8. Display Results
	ui.PrintResults(results)
}