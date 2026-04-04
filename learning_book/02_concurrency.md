# Chapter 2: Concurrency

## Context
A full account search requires fetching data from the Bank, Shared Inventory, Material Storage, and Characters. Doing this sequentially (one after the other) was too slow because each network request takes hundreds of milliseconds.

## Key Go Concepts Learned

### 1. Goroutines
Go makes concurrent programming incredibly easy using "goroutines". Adding the keyword `go` in front of a function call runs it asynchronously in the background.

```go
go s.client.GetSharedInventory()
go s.client.GetBank()
```

### 2. sync.WaitGroup
If you just fire off goroutines, the `main` function will exit before they finish. We learned to use `sync.WaitGroup` to wait for all background tasks to complete.

```go
var wg sync.WaitGroup
wg.Add(4) // We have 4 concurrent tasks

go func() {
    defer wg.Done() // Decrements the counter when the function finishes
    shared, errShared = s.client.GetSharedInventory()
}()
// ... fire off other 3 tasks ...

wg.Wait() // Blocks here until the counter reaches 0
```

### 3. Data Races & Thread Safety
When multiple goroutines write to the same map or slice at the exact same time, it causes a "data race" and the program will crash. We learned two ways to handle this:
1. Return separate variables from each goroutine and merge them safely *after* `wg.Wait()`.
2. Use a `sync.Mutex` to lock the data structure before writing to it.
