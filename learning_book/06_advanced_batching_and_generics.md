# Chapter 6: Advanced Batching & The Interface Trap

## Context
When we implemented "Collection Unlocks" (Phase 17), we needed to resolve thousands of IDs for Skins, Dyes, Minis, and Mounts. To keep the code clean, I initially tried to write a generic helper function called `getEntities`. 

The idea was simple: pass in a URL path and a slice pointer, and let the helper handle the batching and API calls. It worked for the first 200 items, but for anything larger, it only returned the *last* batch. 

## Key Go Concepts Learned

### 1. The Danger of `interface{}` with Slices
In Go, an `interface{}` (or the newer `any`) can hold anything. However, if you pass a slice into a function as an `interface{}`, you lose the ability to easily `append` to it unless you use reflection (which is slow and complex).

In the broken `getEntities` implementation:
```go
func (c *Client) getEntities(path string, ids []int, result interface{}) error {
    // ... batching loop ...
    resp, err := c.rest.R().
        SetQueryParam("ids", strings.Join(batch, ",")).
        SetResult(result). // Trap!
        Get(path)
}
```
The `SetResult(result)` call from the `resty` library was unmarshaling the JSON into the `result` interface. Because this was inside a loop, every new batch was **overwriting** the data from the previous batch instead of appending to it.

### 2. Explicit Batching & Appending
To fix this, we moved away from the "too clever" generic helper and back to explicit type-safe loops. By defining the slice *inside* the calling function and appending to it manually after each batch, we ensure all data is captured.

**The Correct Pattern:**
```go
func (c *Client) ResolveColors(ids []int) ([]NamedEntity, error) {
    var allItems []NamedEntity
    batchSize := 200

    for i := 0; i < len(ids); i += batchSize {
        // ... calculate batch range ...
        var batchItems []NamedEntity
        resp, err := c.rest.R().
            SetResult(&batchItems). // Unmarshal into local batch slice
            Get("/colors")
        
        // Append local batch to the master list
        allItems = append(allItems, batchItems...)
    }
    return allItems, nil
}
```

### 3. Practical Lesson: Don't Abstraction Too Early
This bug was a classic case of "DRY" (Don't Repeat Yourself) backfiring. By trying to save 10 lines of code with a generic helper, I introduced a subtle data-loss bug. In Go, it is often better to have a bit of repetitive, clear, and type-safe code than a single complex, generic function that hides its behavior.
