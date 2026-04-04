# Chapter 1: Basics & API Clients

## Context
When we started GW2CLI, the goal was to fetch data from the Guild Wars 2 API. This required understanding Go's project structure, defining data models, and making HTTP requests.

## Key Go Concepts Learned

### 1. Go Modules
Go uses modules for dependency management. Running `go mod init gw2cli` created a `go.mod` file. This acts as the anchor for the project, tracking dependencies like `github.com/go-resty/resty/v2` which we used for HTTP requests.

### 2. Structs and JSON Tags
APIs return JSON. Go is strongly typed, so we need to map JSON fields to Go structs.
```go
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
```
The `` `json:"id"` `` part is a struct tag. It tells Go's JSON unmarshaler exactly which JSON field maps to which struct property.

### 3. Error Handling
Go doesn't use `try/catch` exceptions. Instead, errors are returned as normal values. This forces you to handle them explicitly.
```go
resp, err := client.Get("/account/inventory")
if err != nil {
    return nil, err // Network error
}
if resp.IsError() {
    return nil, fmt.Errorf("API error: %s", resp.Status()) // HTTP 4xx/5xx error
}
```

### 4. HTTP Clients & Retries
We learned that raw network calls are fragile. We used the `resty` package to abstract standard HTTP calls and easily add automatic retries and exponential backoff, which makes the CLI resilient against rate limits (HTTP 429) or temporary GW2 API outages.
