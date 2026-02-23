# GW2CLI - Guild Wars 2 Inventory Search Tool

A fast and modular CLI tool written in Go to search your entire Guild Wars 2 account. It scans your **Bank**, **Shared Inventory Slots**, **Character Bags**, **Equipped Gear**, and **Account Wallet**.

## Features

- 🔍 **Full Account Search**: Scans Bank, Shared Slots, and all Characters (Bags + Equipment).
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

### Build from Source (Windows)

Using **PowerShell**:
```powershell
git clone https://github.com/fara-moeeni/gw2cli.git
cd gw2cli
go build -o gw2cli.exe cmd/gw2cli/main.go
```

If you have `make` installed (e.g., via Chocolatey or MSYS2), you can also run `make build`.

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

**List all items across your entire account:**
```bash
# Linux/macOS
./gw2cli

# Windows (PowerShell)
.\gw2cli.exe
```

**Search for a specific item (by name, type, or ID):**
```bash
./gw2cli sword
./gw2cli -item "mystic coin"
./gw2cli -item 19976
```

**Filter items by type or character:**
```bash
./gw2cli -item "mystic coin" -character "MyWarrior"
./gw2cli -type Weapon
```

**List all unique item types currently in your inventory:**
```bash
./gw2cli -list-types
```

**Show all assets in your account wallet (Gold, Gems, Karma, etc.):**
```bash
# Linux/macOS
./gw2cli -wallet

# Windows (PowerShell)
.\gw2cli.exe -wallet
```

**List all characters on your account with details:**
```bash
./gw2cli -list-characters
```

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
