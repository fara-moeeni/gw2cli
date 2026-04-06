# GW2CLI - Guild Wars 2 Inventory Search Tool

A fast and modular CLI tool written in Go to search your entire Guild Wars 2 account. It scans your **Bank**, **Shared Inventory Slots**, **Material Storage**, **Character Bags**, **Equipped Gear**, and **Account Wallet**.

## Features

- 🔍 **Full Account Search**: Scans Bank, Shared Slots, Material Storage, and all Characters (Bags + Equipment).
- 👤 **Character Management**: List all characters with level, profession, and playtime details.
- 💰 **Wallet View**: Quickly view all account currencies and their values.
- 🚀 **Fast**: Uses concurrent requests to fetch account data and parallel processing for item resolution.
- 📦 **Smart Batching**: Automatically handles API limits when resolving thousands of item IDs.
- 📝 **Flexible Filtering**: Search by Name, Item Type, ID, or Location.

## Installation

### Prerequisites
- [Go 1.23+](https://go.dev/dl/)

### Build from Source (Linux/macOS)

```bash
git clone https://github.com/fara-moeeni/gw2cli.git
cd gw2cli
make build
```

### Manual (man) Pages (Linux/macOS)

You can generate and install a `man` page for `gw2cli` to access offline documentation.

**Generate the man page (requires `go-md2man`):**
```bash
# Install tool if missing
go install github.com/cpuguy83/go-md2man/v2@latest

# Generate and install (requires sudo for /usr/local/share/man/man1/)
make man
sudo make install
```

### Build from Source (Windows)

Using **PowerShell**:
```powershell
git clone https://github.com/fara-moeeni/gw2cli.git
cd gw2cli
go build -o gw2cli.exe cmd/gw2cli/main.go
```

If you have `make` installed (e.g., via Chocolatey or MSYS2), you can also run `make build`.

### Triggering a Release (Maintainers Only)

The GitHub Actions pipeline automatically builds and cross-compiles binaries when a new version tag is pushed:

```bash
git checkout main
git pull origin main
git tag 2.x.x
git push origin 2.x.x
```

---

## Usage

### 1. Get an API Key
You need a Guild Wars 2 API Key with at least the following permissions:
- `account`
- `inventories`
- `characters`
- `wallet`

Generate one at: [https://account.arena.net/applications](https://account.arena.net/applications)

### 2. Set Environment Variable
Export your API key as an environment variable:

**Linux/macOS (Bash/Zsh):**
```bash
export GW2_API_KEY="YOUR_API_KEY_HERE"
man gw2cli # View full documentation (if installed)
```

**Windows (PowerShell):**
```powershell
$env:GW2_API_KEY="YOUR_API_KEY_HERE"
```

**Windows (Command Prompt):**
```cmd
set GW2_API_KEY=YOUR_API_KEY_HERE
```

### 3. Examples

Below are examples for **Linux/macOS**. For **Windows**, replace `./gw2cli` with `.\gw2cli.exe` in PowerShell or `gw2cli.exe` in Command Prompt.

**Search for an item across your entire account:**
```bash
./gw2cli search "mystic coin"
./gw2cli search sword -type Weapon
./gw2cli search -character "MyWarrior"
```

**List account information:**
```bash
./gw2cli list characters
./gw2cli list types
```

**Show account wallet (Gold, Gems, Karma, etc.):**
```bash
./gw2cli wallet
```

**Check Trading Post delivery, orders, and history:**
```bash
./gw2cli tp delivery
./gw2cli tp orders
./gw2cli tp history
```

**Check Trading Post price for an item:**
```bash
./gw2cli tp price "Mystic Coin"
```

**Check current Gem/Coin exchange rates:**
```bash
./gw2cli exchange
./gw2cli exchange gems 400
./gw2cli exchange coins 1000000
```

**Check Legendary Armory:**
```bash
./gw2cli legendary
./gw2cli legendary "Twilight"
```

**Build and search the local item database:**
```bash
# Build/Refresh cache
./gw2cli cache update

# Offline name-based search
./gw2cli cache find "Twilight"
```

**Search unlocked crafting recipes:**
```bash
# List all unlocked recipes
./gw2cli recipes

# Search unlocked recipes by output item name
./gw2cli recipes find "Deldrimor"

# Find unlocked recipes that use a specific ingredient
./gw2cli recipes ingredient "Ectoplasm"
```

**Check account collections:**
```bash
# Show summary of all collections
./gw2cli collection

# List all unlocked minis
./gw2cli collection minis

# Search for "Dragon" skins
./gw2cli collection skins "Dragon"

# List unlocked dyes
./gw2cli collection dyes
```

**Track daily and weekly activities:**
```bash
# Show daily summary (Bosses, Dungeons, Fractals, Wizard's Vault)
./gw2cli daily

# Show only today's fractal dailies
./gw2cli daily fractals

# Show only world boss daily reset status
./gw2cli daily bosses

# Show daily Wizard's Vault objectives
./gw2cli daily wizardsvault

# Show weekly summary (Raids, Wizard's Vault weekly)
./gw2cli weekly
```
Note: Convergences rotation is not available via the API — check [https://wiki.guildwars2.com/wiki/Convergences](https://wiki.guildwars2.com/wiki/Convergences)

---

## Project Structure

The project follows a modular Go structure:

```
.
├── cmd/
│   └── gw2cli/
│       └── main.go       # CLI Entry point & flag parsing
├── internal/
│   ├── inventory/        # Core business logic: search, aggregation, wallet
│   └── ui/               # Formatted output and help systems
├── pkg/
│   └── gw2api/           # Guild Wars 2 API client implementation
├── go.mod                # Go module definition
└── Makefile              # Build automation
```

## License

MIT
