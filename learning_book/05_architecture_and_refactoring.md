# Chapter 5: Architecture & Refactoring

## Context
As GW2CLI grew to include searching, wallets, Trading Post queries, gem exchanges, and cache management, the `main.go` file became a massive `if-else` chain. Furthermore, the command syntax became confusing (e.g., `gw2cli "sword" -type Weapon` vs `gw2cli -tp-price "sword"`). We needed a standard "Subcommand" architecture.

## Key Go Concepts Learned

### 1. Subcommands via `flag.NewFlagSet`
Instead of one global `flag.Parse()`, Go allows you to create independent sets of flags for different commands (like `git commit` vs `git push`).

```go
// Define a flagset
searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
searchType := searchCmd.String("type", "", "Filter by strict Item Type")

// Define another flagset
walletCmd := flag.NewFlagSet("wallet", flag.ExitOnError)
```

### 2. Command Routing
We learned to read `os.Args` directly to figure out which subcommand the user wants to run, and then pass the *rest* of the arguments to that specific FlagSet.

```go
switch os.Args[1] {
case "search":
    searchCmd.Parse(os.Args[2:]) // Parse only the arguments after 'search'
    // Execute search logic
case "wallet":
    walletCmd.Parse(os.Args[2:])
    // Execute wallet logic
}
```

### 3. The Value of Refactoring
This refactor (Phase 13) was a major "breaking change" (version bumped to v2.0.0). However, the business logic inside `internal/inventory` barely changed.
*Lesson Learned: By keeping our core business logic separated in `internal/inventory` and cleanly decoupled from the CLI layer (`cmd/gw2cli/main.go`), we were able to completely rewrite the user interface without breaking how the app actually interacts with the GW2 API.*

### 4. Custom Help Menus
With FlagSets, we could assign specific help text to specific commands by overriding the `.Usage` function:
```go
searchCmd.Usage = ui.PrintSearchHelp
```
This made the CLI much friendlier, as users no longer had to read through irrelevant Trading Post flags when they just wanted help with the `search` command.
