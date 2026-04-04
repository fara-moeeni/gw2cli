# Chapter 4: File I/O & Caching

## Context
Fetching Trading Post prices by name or resolving 70,000 item IDs requires downloading the entire item database. Doing this on every run was impossible. We needed a local JSON cache stored on the user's filesystem.

## Key Go Concepts Learned

### 1. The `os` and `filepath` Packages
We learned how to interact securely with the operating system to find the user's home directory and create cross-platform paths.
```go
home, err := os.UserHomeDir()
cachePath := filepath.Join(home, ".config", "gw2cli", "items.json")
```

### 2. Creating Directories and Files
Before writing a file, the parent directories must exist.
```go
// 0755 permissions: User can read/write/execute, others can read/execute
os.MkdirAll(filepath.Dir(cachePath), 0755)

// 0644 permissions: User can read/write, others can read
os.WriteFile(cachePath, jsonData, 0644)
```

### 3. Cache Expiry with `time`
We learned how to check file metadata to warn the user if the cache is stale.
```go
info, err := os.Stat(cachePath)
// time.Since calculates duration from a point in the past until now
if time.Since(info.ModTime()) > 7*24*time.Hour {
    fmt.Println("warning: item cache is 7+ days old")
}
```

### 4. JSON Serialization (`json.Marshal`)
To save Go structs to disk, we convert them to bytes using `json.Marshal`. To load them back into memory, we use `json.Unmarshal`. We kept the structs minimal (ID, Name, Type) to save disk space and reduce memory footprint when loading the 70k+ items into RAM.
