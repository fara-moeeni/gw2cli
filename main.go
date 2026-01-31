package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"

	"gw2cli/pkg/gw2"
)

type ItemLocation struct {
	Location string
	Count    int
}

func printGlobalHelp() {
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
        Show this help message.

Examples:
  ./gw2cli "mystic coin"                  # Global fuzzy search
  ./gw2cli -item "Sword" -type "Weapon"   # Find weapons named "Sword"
  ./gw2cli -character "Bank"              # List everything in the Bank
  ./gw2cli -character "Rytlock" -type "Armor" # Find Armor on Rytlock`)
}

func printTypeHelp() {
	fmt.Println(`
Option: -type [Category]

Description:
  Filters search results to strictly match a specific Item Type.
  This is useful for distinguishing between an item's Name and its Category.
  
  For example, searching for "Sword" might return a "Wood Sword" (Weapon) 
  and "Sword Oil" (Consumable). Using -type "Weapon" ensures you only see weapons.

Usage:
  ./gw2cli -type "Weapon"
  ./gw2cli -type "Armor" -item "Boots"

Available Types (common examples):
  - Armor
  - Weapon
  - Consumable
  - UpgradeComponent
  - CraftingMaterial
  - Trinket
  - Bag
  - Container
  - Trophy

To see ALL types currently in your inventory, run:
  ./gw2cli -list-types`)
}

func printItemHelp() {
	fmt.Println(`
Option: -item [SearchTerm]

Description:
  Search for an item by its Name or ID.
  This is a case-insensitive "fuzzy" search.

Usage:
  ./gw2cli -item "Berserker"     # Matches "Berserker's Helm", "Berserker's Sword"
  ./gw2cli -item 12345           # Matches Item ID 12345
  ./gw2cli -item "Mystic Coin"   # Specific name match`)
}

func printCharacterHelp() {
	fmt.Println(`
Option: -character [Name/Location]

Description:
  Narrow down the search to a specific character, storage location, or bag.
  Matches against the "Location" field of the item.

Usage:
  ./gw2cli -character "Rytlock"       # Only items on character "Rytlock"
  ./gw2cli -character "Bank"          # Only items in the Bank
  ./gw2cli -character "Shared"        # Only items in Shared Inventory Slots
  ./gw2cli -character "Equipped"      # Only items currently equipped (on any character)`)
}

func main() {
	// Custom Usage function for standard -h / --help
	flag.Usage = printGlobalHelp

	// 1. Parse Flags immediately
	itemFlag := flag.String("item", "", "Search by Item Name or ID (use '-item help' for details)")
	typeFlag := flag.String("type", "", "Filter by strict Item Type (use '-type help' for details)")
	charFlag := flag.String("character", "", "Search by Character Name or Location (use '-character help' for details)")
	listTypesFlag := flag.Bool("list-types", false, "List all unique item types found in your inventory")
	
	// Explicitly handle a custom boolean flag if user prefers explicit -help
	// (Though flag.Parse handles -help automatically if we don't define it, defining it lets us control the output perfectly)
	helpFlag := flag.Bool("help", false, "Show this help message")

	flag.Parse()

	// 2. Handle Help Requests
	if *helpFlag {
		printGlobalHelp()
		return
	}
	if strings.ToLower(*typeFlag) == "help" {
		printTypeHelp()
		return
	}
	if strings.ToLower(*itemFlag) == "help" {
		printItemHelp()
		return
	}
	if strings.ToLower(*charFlag) == "help" {
		printCharacterHelp()
		return
	}

	// Positional args act as a "global" filter
	globalFilter := strings.ToLower(strings.Join(flag.Args(), " "))

	apiKey := os.Getenv("GW2_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set the GW2_API_KEY environment variable")
	}

	client := gw2.NewClient(apiKey)

	// Fetch everything in parallel
	var wg sync.WaitGroup
	var shared gw2.AccountInventory
	var bank gw2.AccountInventory
	var characters []gw2.Character
	var errShared, errBank, errChars error

	wg.Add(3)

	fmt.Println("Fetching account data...")

	go func() {
		defer wg.Done()
		shared, errShared = client.GetSharedInventory()
	}()

	go func() {
		defer wg.Done()
		bank, errBank = client.GetBank()
	}()

	go func() {
		defer wg.Done()
		characters, errChars = client.GetCharacters()
	}()

	wg.Wait()

	if errShared != nil {
		log.Printf("Warning: Could not fetch shared inventory: %v", errShared)
	}
	if errBank != nil {
		log.Printf("Warning: Could not fetch bank: %v", errBank)
	}
	if errChars != nil {
		log.Printf("Warning: Could not fetch characters: %v", errChars)
	}

	// Map ItemID -> List of Locations
	inventoryMap := make(map[int][]ItemLocation)
	var allIDs []int

	// Process Shared
	for i, slot := range shared {
		if slot != nil && slot.ID != 0 {
			allIDs = append(allIDs, slot.ID)
			inventoryMap[slot.ID] = append(inventoryMap[slot.ID], ItemLocation{
				Location: fmt.Sprintf("Shared Slot %d", i+1),
				Count:    slot.Count,
			})
		}
	}

	// Process Bank
	for i, slot := range bank {
		if slot != nil && slot.ID != 0 {
			allIDs = append(allIDs, slot.ID)
			inventoryMap[slot.ID] = append(inventoryMap[slot.ID], ItemLocation{
				Location: fmt.Sprintf("Bank Tab %d", (i/30)+1),
				Count:    slot.Count,
			})
		}
	}

	// Process Characters
	for _, char := range characters {
		// Bags
		for bagIdx, bag := range char.Bags {
			if bag != nil {
				for _, slot := range bag.Inventory {
					if slot != nil && slot.ID != 0 {
						allIDs = append(allIDs, slot.ID)
						inventoryMap[slot.ID] = append(inventoryMap[slot.ID], ItemLocation{
							Location: fmt.Sprintf("%s (Bag %d)", char.Name, bagIdx+1),
							Count:    slot.Count,
						})
					}
				}
			}
		}
		// Equipment
		for _, equip := range char.Equipment {
			if equip.ID != 0 {
				allIDs = append(allIDs, equip.ID)
				inventoryMap[equip.ID] = append(inventoryMap[equip.ID], ItemLocation{
					Location: fmt.Sprintf("%s (Equipped: %s)", char.Name, equip.Slot),
					Count:    1,
				})
			}
		}
	}

	if len(allIDs) == 0 {
		fmt.Println("No items found on account.")
		return
	}

	fmt.Printf("Resolving %d items details...\n", len(allIDs))
	items, err := client.GetItems(allIDs)
	if err != nil {
		log.Fatalf("Error resolving items: %v", err)
	}

	// Handle -list-types flag
	if *listTypesFlag {
	
uniqueTypes := make(map[string]bool)
		for _, item := range items {
		
uniqueTypes[item.Type] = true
		}

		var sortedTypes []string
		for t := range uniqueTypes {
			sortedTypes = append(sortedTypes, t)
		}
		sort.Strings(sortedTypes)

		fmt.Println("\n--- Available Item Types (in your inventory) ---")
		for _, t := range sortedTypes {
			fmt.Println(t)
		}
		return // Exit after listing types
	}

	fmt.Println("\n--- Search Results ---")
	found := false

	// Normalize flags
	targetItem := strings.ToLower(*itemFlag)
	targetType := strings.ToLower(*typeFlag)
	targetChar := strings.ToLower(*charFlag)

	for _, item := range items {
		// 1. Strict Type Filter
		if targetType != "" {
			if strings.ToLower(item.Type) != targetType {
				continue
			}
		}

		// 2. Item Name/ID Filter
		if targetItem != "" {
			// Check Name OR ID (Type is handled separately now, unless user didn't use -type flag)
			match := strings.Contains(strings.ToLower(item.Name), targetItem) ||
				strings.Contains(fmt.Sprintf("%d", item.ID), targetItem)
			
			// Fallback: If they didn't specify -type, allow -item to match Type too (legacy behavior)
			if targetType == "" && strings.Contains(strings.ToLower(item.Type), targetItem) {
				match = true
			}

			if !match {
				continue
			}
		}

		// 3. Location Level Filter
		var validLocations []ItemLocation
		for _, loc := range inventoryMap[item.ID] {
			if targetChar != "" {
				if !strings.Contains(strings.ToLower(loc.Location), targetChar) {
					continue
				}
			}
			validLocations = append(validLocations, loc)
		}

		if len(validLocations) == 0 {
			continue
		}

		// 4. Global Filter (Positional Args)
		var locStrings []string
		totalCount := 0
		for _, loc := range validLocations {
			locStrings = append(locStrings, fmt.Sprintf("%s (x%d)", loc.Location, loc.Count))
			totalCount += loc.Count
		}
		
		fullText := fmt.Sprintf("%s %s %d %s", item.Name, item.Type, item.ID, strings.Join(locStrings, " "))

		if globalFilter == "" || strings.Contains(strings.ToLower(fullText), globalFilter) {
			found = true
			fmt.Printf("\nName: \033[1m%s\033[0m (ID: %d)\n", item.Name, item.ID)
			fmt.Printf("Type: %s\n", item.Type)
			fmt.Printf("Total: %d\n", totalCount)
			fmt.Println("Found in:")
			for _, loc := range validLocations {
				fmt.Printf(" - %s: %d\n", loc.Location, loc.Count)
			}
		}
	}

	if !found {
		fmt.Println("No items found matching your criteria.")
	}
}