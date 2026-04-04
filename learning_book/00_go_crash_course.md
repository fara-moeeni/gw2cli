# Chapter 0: Go Crash Course

## Context
Before we dive into building GW2CLI, it's essential to understand the basic syntax and philosophy of the Go programming language. Go (or Golang) was designed by Google to be simple, fast, and highly readable. It is a statically typed, compiled language.

## Basic Syntax

### 1. Variables and Types
In Go, variables are explicitly typed, but the compiler can also infer the type.
```go
// Explicitly declaring a variable
var name string = "GW2CLI"
var version int = 2

// Short variable declaration (infers type, only works inside functions)
author := "fara-moeeni"
isAwesome := true
```

### 2. Functions
Functions in Go can return multiple values, which is heavily used for error handling.
```go
// A simple function
func add(x int, y int) int {
    return x + y
}

// A function returning multiple values (common for returning a result and an error)
func divide(x int, y int) (int, error) {
    if y == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return x / y, nil
}
```

### 3. Control Structures
Go simplifies control structures. There is no `while` loop; `for` does everything.

**If Statements:**
```go
if version >= 2 {
    fmt.Println("Using subcommand architecture")
} else {
    fmt.Println("Using legacy flat flags")
}
```

**For Loops:**
```go
// Traditional for loop
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// "While" loop equivalent
x := 0
for x < 5 {
    x++
}

// Ranging over a collection (like a list of items)
items := []string{"Sword", "Shield", "Staff"}
for index, item := range items {
    fmt.Printf("Item %d: %s\n", index, item)
}
```

### 4. Collections: Arrays, Slices, and Maps
- **Arrays** have a fixed size. We rarely use them directly.
- **Slices** are dynamic, resizable arrays. This is what we use 99% of the time (e.g., `[]string`).
- **Maps** are key-value stores (like dictionaries in Python or objects in JavaScript).

```go
// Slice (dynamic array)
inventory := []string{"Potion", "Food"}
inventory = append(inventory, "Booster") // Adding an item

// Map (key-value store)
itemPrices := make(map[string]int)
itemPrices["Mystic Coin"] = 20000 // 20000 copper (2 gold)
fmt.Println(itemPrices["Mystic Coin"])
```

### 5. Structs and Methods (Object-Oriented Go)
Go doesn't have "classes". Instead, it has `structs` (custom data types combining different fields) and you can attach functions to them, known as `methods`.

```go
// Define a struct
type Character struct {
    Name  string
    Level int
}

// Define a method attached to the Character struct
func (c Character) LevelUp() {
    c.Level++
    fmt.Println(c.Name, "is now level", c.Level)
}
```
*Note: We will learn about Pointers in Chapter 0.1 to understand how to actually modify the struct's original data in a method.*
