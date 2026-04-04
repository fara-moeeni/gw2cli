# Chapter 3: CLI Flags & Formatting

## Context
To interact with the CLI, we needed a way to accept user input. We started with Go's standard `flag` package, creating a flat structure where every feature had a global flag (e.g., `-wallet`, `-item`, `-list-characters`).

## Key Go Concepts Learned

### 1. The `flag` Package
Go's `flag` package binds command-line arguments to pointers.
```go
walletFlag := flag.Bool("wallet", false, "Show account wallet")
itemFlag := flag.String("item", "", "Search by Item Name")
flag.Parse() // This actually reads os.Args and populates the pointers
```
Because they are pointers, you have to dereference them (using `*walletFlag`) to get the actual value.

### 2. Edge Cases in Flat Flags
We quickly learned that a flat flag structure is brittle. If a user typed `gw2cli -item -wallet`, the `flag` package would consume `-wallet` as the string value for `-item`, leading to confusing bugs.
To fix this, we had to write complex edge-case logic:
```go
if strings.HasPrefix(*itemFlag, "-") {
    // Manually intercept flags that were accidentally parsed as values
}
```
*Lesson Learned: Flat flag structures don't scale well for complex CLIs.* (See Chapter 5).

### 3. Console Formatting
We learned how to make output look professional using `fmt.Printf` and formatting verbs.
```go
// %-20s pads the string to 20 characters, aligning columns perfectly
fmt.Printf("% -20s %-15s\n", char.Name, char.Profession)
```
We also learned basic terminal color codes (ANSI escape sequences) to make text bold: `\033[1m%s\033[0m`.
