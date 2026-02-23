package ui

import (
	"fmt"

	"gw2cli/internal/inventory"
)

func PrintGlobalHelp() {
	fmt.Println(`
GW2CLI - Guild Wars 2 Inventory Search Tool (v1.2.0)

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

  -list-characters
        List all characters on your account with details (Level, Profession, etc).

  -wallet
        Show all currencies in your account wallet.

  -tp-delivery
        Show items and coin waiting to be picked up from the Trading Post.

  -tp-orders
        Show current active buy and sell orders.

  -tp-history
        Show past transaction history (bought and sold items).

  -tp-price <item name or ID>
        Look up the current buy and sell price for a specific item.

  -help
        Show this help message.`)
}

func PrintTypeHelp() {
	fmt.Println(`
Option: -type [Category]

Description:
  Filters search results to strictly match a specific Item Type.
  `)
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
