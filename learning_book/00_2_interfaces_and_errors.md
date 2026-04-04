# Chapter 0.2: Interfaces & Errors

## Context
In Go, errors are just values that implement a specific interface. Understanding interfaces is crucial to understanding how Go handles polymorphism and error management.

## Key Concepts

### 1. What is an Interface?
An interface is a type that specifies a set of method signatures. If a struct implements all the methods defined in an interface, it automatically satisfies that interface. There is no `implements` keyword in Go!

```go
// Define an interface
type Describer interface {
    Describe() string
}

// A struct
type Item struct {
    Name string
}

// Item implements the Describer interface because it has a Describe() string method
func (i Item) Describe() string {
    return "Item: " + i.Name
}

func printDescription(d Describer) {
    fmt.Println(d.Describe())
}

func main() {
    mySword := Item{Name: "Twilight"}
    printDescription(mySword) // Works!
}
```

### 2. The `error` Interface
The most common interface you will encounter in Go is the built-in `error` interface. It looks like this under the hood:
```go
type error interface {
    Error() string
}
```
This means that *any* type that has an `Error() string` method can be treated as an error! 

### 3. Handling Errors in GW2CLI
When we make an API call or read a file, we often see this pattern:
```go
data, err := os.ReadFile("items.json")
if err != nil {
    // Handle the error
    fmt.Println("Failed to read file:", err.Error())
    return
}
```
Because `err` is an interface, it could be a simple string error, or a complex struct containing detailed information about a network timeout. As long as it has the `Error()` method, Go treats it as an error.

### 4. Creating Custom Errors
We can use the `fmt` package to create formatted errors easily, which we do frequently in GW2CLI:
```go
if quantity < 0 {
    return fmt.Errorf("invalid quantity: %d", quantity)
}
```
