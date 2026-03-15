# gw2cli(1) | GW2CLI User Manual

## NAME
gw2cli - A command-line interface for the Guild Wars 2 API.

## SYNOPSIS
**gw2cli** [FLAGS] [SEARCH_TERM]

## DESCRIPTION
**gw2cli** provides a fast, terminal-based interface to manage your Guild Wars 2 account. It supports searching through your bank, shared inventory, and character inventories, as well as checking Trading Post prices, wallet balances, and currency exchange rates.

The tool uses a local item cache to provide fast offline searching and to resolve item names to IDs.

## OPTIONS

### General Flags
**-help**
    Show the help message and exit.

**-verbose**
    Enable verbose output for debugging and detailed progress tracking.

**-update-cache**
    Fetch and update the local item database from the GW2 API. Recommended to run once a week.

**-find** *TERM*
    Search for items in the local database by name without calling the API.

### Account & Character Flags
**-list-characters**
    List all characters on the account with their level, profession, and total playtime.

**-wallet**
    Show the account's wallet (Gold, Gems, Karma, etc.).

**-character** *NAME*
    Filter search results to a specific character.

### Inventory & Search
**-item** *TERM*
    Search for an item by name or ID across all account inventories.

**-type** *CATEGORY*
    Filter search results by a strict item type (e.g., "Weapon", "Armor").

**-list-types**
    List all unique item types currently found in your inventories.

### Trading Post (Commerce)
**-tp-price** *NAME_OR_ID*
    Check current buy and sell offers on the Trading Post.

**-tp-delivery**
    Show items and coins waiting to be picked up from the Trading Post.

**-tp-orders**
    Show active buy and sell orders.

**-tp-history**
    Show the last 90 days of Trading Post transaction history.

### Currency Exchange
**-exchange**
    Show current Gem-to-Gold and Gold-to-Gem exchange rates.

**-exchange-gems** *AMOUNT*
    Calculate how much gold you would receive for a specific amount of gems.

**-exchange-coins** *AMOUNT*
    Calculate how many gems you would receive for a specific amount of gold (in copper).

## EXAMPLES
**Search for "Legendary" items across the account:**
    gw2cli "Legendary"

**Check the price of "Ectoplasm":**
    gw2cli -tp-price "Glob of Ectoplasm"

**List all characters and their playtime:**
    gw2cli -list-characters

**Calculate exchange for 400 gems:**
    gw2cli -exchange-gems 400

## FILES
**~/.config/gw2cli/items.json**
    The local item cache used for name resolution and offline search.

## ENVIRONMENT
**GW2_API_KEY**
    The API key used to authenticate with the Guild Wars 2 API. Must have 'account', 'inventories', 'characters', 'tradingpost', and 'wallet' permissions.

## AUTHOR
fara-moeeni (https://github.com/fara-moeeni/gw2cli)
