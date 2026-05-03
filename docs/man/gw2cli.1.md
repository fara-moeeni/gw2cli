# gw2cli(1) | GW2CLI User Manual

## NAME
gw2cli - A command-line interface for the Guild Wars 2 API.

## SYNOPSIS
**gw2cli** <COMMAND> [ARGUMENTS] [FLAGS]

## DESCRIPTION
**gw2cli** provides a fast, terminal-based interface to manage your Guild Wars 2 account. It uses a subcommand architecture to organize its various features, including inventory search, character listing, Trading Post management, and currency exchange.

The tool uses a local item cache to provide fast offline searching and to resolve item names to IDs.

## GLOBAL OPTIONS
**-verbose**
    Enable verbose output for debugging and detailed progress tracking.

**-help**
    Show the global help message and exit.

**version**
    Show the version information and exit.

## COMMANDS

### account
Show account summary, including fractal level, WvW rank, AP, account age, and playtime.

### search [TERM] [FLAGS]
Search your entire account inventory (Bank, Shared Inventory, and all Characters) for items.

**-type** *CATEGORY*
    Filter search results by a strict item type (e.g., "Weapon", "Armor").

**-character** *NAME*
    Filter search results to a specific character.

### list [types|characters|account]
List high-level account information.

**types**
    List all unique item types currently found in your inventories.

**characters**
    List all characters on the account with their level, profession, and total playtime.

**account**
    Show the account summary.

### wallet
Show the account's wallet (Gold, Gems, Karma, etc.).

### tp <subcommand> [ARGUMENTS]
Manage Trading Post (Commerce) activities.

**delivery**
    Show items and coins waiting to be picked up from the Trading Post.

**orders**
    Show current active buy and sell orders.

**history**
    Show the last 90 days of Trading Post transaction history.

**price** *NAME_OR_ID*
    Check current buy and sell offers on the Trading Post for a specific item.

### exchange [subcommand] [AMOUNT]
Check Gem/Coin exchange rates.

**gems** *AMOUNT*
    Calculate how many coins you would receive for a specific amount of gems.

**coins** *AMOUNT*
    Calculate how many gems you would receive for a specific amount of coins (in copper).

### cache <subcommand> [ARGUMENTS]
Manage the local item database.

**update**
    Fetch and update the local item database from the GW2 API. Recommended to run once a week.

**find** *TERM*
    Search for items in the local database by name without calling the API.

### achievements [subcommand] [FLAGS]
Track achievement progress and masteries.

**all**
    List all achievements started or completed on the account.
    **--status**=[completed|incomplete|any]
        Filter the list by completion status. Default is "any".

**find** *TERM*
    Search achievements by name and show progress.

**update-cache**
    Build or refresh the local achievement database.

**masteries**
    Show mastery points earned across all regions and total luck.

### daily [subcommand]
Track daily reset status and objectives.

**bosses**
    List world bosses. With an API key, completed bosses are marked.

**fractals**
    Show current daily fractal achievements.

**wizardsvault**
    Show daily Wizard's Vault objectives. Requires an API key.

### weekly
Show weekly raid and Wizard's Vault status. Requires an API key.

### recipes [subcommand] [ARGUMENTS]
List and search unlocked account recipes. Requires an API key and local item cache.

### collection [subcommand] [FILTER]
Show account collection summaries and unlocked skins, dyes, minis, mounts, outfits, novelties, or finishers. Requires an API key.

### legendary [TERM]
Display all items currently in your Legendary Armory. Requires an API key and local item cache.

## UNAUTHENTICATED COMMANDS
The following commands work without `GW2_API_KEY`: `exchange`, `tp price`, `daily bosses`, `daily fractals`, `cache update`, `cache find`, and `achievements update-cache`.

## EXAMPLES
**Search for "Legendary" items across the account:**
    gw2cli search "Legendary"

**Check the price of "Ectoplasm":**
    gw2cli tp price "Glob of Ectoplasm"

**List all characters and their playtime:**
    gw2cli list characters

**Calculate exchange for 400 gems:**
    gw2cli exchange gems 400

## FILES
**~/.config/gw2cli/items.json**
    The local item cache used for name resolution and offline search.

## ENVIRONMENT
**GW2_API_KEY**
    The API key used to authenticate with the Guild Wars 2 API. Must have 'account', 'inventories', 'characters', 'tradingpost', and 'wallet' permissions.

## AUTHOR
fara-moeeni (https://github.com/fara-moeeni/gw2cli)
