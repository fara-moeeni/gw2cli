# GW2CLI - Guild Wars 2 Inventory Search Tool

A blazing fast CLI tool written in Go to search your entire Guild Wars 2 account for items. It scans your **Bank**, **Shared Inventory Slots**, **Character Bags**, and **Equipped Gear** to tell you exactly what you have and where it is.

## Features

- 🔍 **Full Account Search**: Scans Bank, Shared Slots, and all Characters (Bags + Equipment).
- 🚀 **Fast**: Uses concurrent requests to fetch account data and parallel processing for item resolution.
- 📦 **Smart Batching**: Automatically handles API limits when resolving thousands of item IDs.
- 📝 **Grep-like Search**: Fuzzy search by Name, Item Type, ID, or Location.

## Installation

### Prerequisites
- [Go 1.23+](https://go.dev/dl/)

### Build from Source

```bash
git clone https://github.com/fara-moeeni/gw2cli.git
cd gw2cli
make build
```

## Release & Versioning

This project follows **Conventional Commits** and **Semantic Versioning**. Releases are automated via GitHub Actions.

### How to Release
1. Commit your changes using conventional messages (e.g., `feat: ...`, `fix: ...`).
2. Tag the commit with a version:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   ```
3. Push the tag:
   ```bash
   git push origin v1.0.0
   ```
This will automatically trigger a GitHub Release with compiled binaries for Linux, macOS, and Windows.

## Usage

### 1. Get an API Key
You need a Guild Wars 2 API Key with the following permissions:
- `account`
- `inventories`
- `characters`

Generate one at: [https://account.arena.net/applications](https://account.arena.net/applications)

### 2. Set Environment Variable
Export your API key as an environment variable:

**Linux/macOS:**
```bash
export GW2_API_KEY="YOUR_API_KEY_HERE"
```

**Windows (PowerShell):**
```powershell
$env:GW2_API_KEY="YOUR_API_KEY_HERE"
```

### 3. Run Commands

**List everything:**
```bash
./gw2cli
```

**Search for an item (by name, type, or ID):**
```bash
./gw2cli sword
./gw2cli "mystic coin"
./gw2cli 12345
```

**Example Output:**
```text
Name: Mystic Coin (ID: 19976)
Type: CraftingMaterial
Total: 250
Found in:
 - Shared Slot 1: 250

Name: Zojja's Claymore (ID: 46774)
Type: Weapon
Total: 1
Found in:
 - MyWarrior (Equipped: WeaponA1): 1
```

## Project Structure

```
.
├── main.go           # CLI Entry point and logic
├── go.mod            # Go module definition
├── pkg/
│   └── gw2/
│       ├── client.go # API Client (Resty)
│       └── types.go  # JSON Data Structures
```

## License

MIT
