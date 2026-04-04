# Chapter 0.1: Pointers & References

## Context
When you pass variables to functions in Go, they are copied by default (passed by value). If you want a function to modify the original variable, or if you want to avoid copying a massive data structure, you must use **Pointers**.

## Key Concepts

### 1. What is a Pointer?
A pointer holds the **memory address** of a value, rather than the value itself.
- `&` generates a pointer to its operand (gives you the address).
- `*` dereferences a pointer (gives you the underlying value).

```go
func main() {
    x := 10
    p := &x         // p is a pointer to x (type *int)
    
    fmt.Println(p)  // Prints a memory address like 0xc00001a0a8
    fmt.Println(*p) // Prints 10 (dereferences the pointer)
    
    *p = 21         // Changes the value at the memory address
    fmt.Println(x)  // Prints 21, because x was modified through the pointer!
}
```

### 2. Pointers in Struct Methods
In the previous chapter, our `LevelUp` method didn't actually permanently modify the character's level because it modified a *copy* of the struct. To modify the original, we use a pointer receiver.

```go
type Character struct {
    Name  string
    Level int
}

// Notice the * before Character. This means we are receiving a pointer.
func (c *Character) LevelUp() {
    c.Level++ 
}

func main() {
    char := Character{Name: "Rytlock", Level: 80}
    char.LevelUp() // This will permanently change the level to 81
}
```

### 3. How Pointers are used in GW2CLI
In GW2CLI, we use pointers extensively in two main areas:
1. **CLI Flags:** The `flag` package returns pointers. `flag.Bool("verbose", false, "...")` returns a `*bool`. When the user runs the CLI, the package updates the value at that memory address. We check it using `*verbose`.
2. **API Responses:** When we fetch data from the Guild Wars 2 API, we pass a pointer to our structs so the JSON library can write the data directly into our variables: `json.Unmarshal(data, &myStruct)`.
