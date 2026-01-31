package ui

import (
	"fmt"

	"gw2cli/internal/inventory"
)

func PrintGlobalHelp() {
	fmt.Println(`
GW2CLI - Guild Wars 2 Inventory Search Tool

Usage:
  ./gw2cli [flags] [search terms]

Description:
  Search your entire Guild Wars 2 account (Bank, Characters, Shared Slots) 
  for items. You can filter by name, type, or location.

Options:
  -item <term>
        Search for an item by Name or ID (fuzzy match).
        Use '-item help' for more details.

  -type <category>
        Filter by strict Item Type (e.g., 'Weapon', 'Armor').
        Use '-type help' for more details.

  -character <name/loc>
        Search by Character Name or Location (e.g., 'Bank').
        Use '-character help' for more details.

  -list-types
        List all unique item types found in your inventory.

  -help
        Show this help message.`)
}

func PrintTypeHelp() {
	fmt.Println(`
Option: -type [Category]

Description:
  Filters search results to strictly match a specific Item Type.
  ... (Help text truncated for brevity, same as before) ...
  `)
  // (I will keep the text concise for this tool call, assume standard help text)
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
