# Chapter 8: Data Merging & Multi-Caching

## Context
By Phase 19, GW2CLI had two local caches: `items.json` and `achievements.json`. The "Achievements" feature required us to merge three different layers of data to show a single result:
1. **Local Cache**: Names and Descriptions (Static, offline).
2. **Account State**: Progress, `Done` status, and current counts (Live, authenticated).
3. **API Metadata**: Achievement Categories and Groups (Public, live).

## Key Go Concepts Learned

### 1. The Multi-Cache Pattern
We used the same XDG-compliant path pattern for achievements as we did for items. This allowed us to reuse the "Staleness Check" logic:
```go
// Check if cache exists and is < 7 days old
if err := invService.CheckAchievementCacheStatus(); err != nil {
    // Prompt to run update-cache
}
```
*Lesson Learned: By standardizing our cache format and logic, we made it easy to add new data types (like Achievements or Masteries) without reinventing the wheel.*

### 2. Merging Data with Maps
To merge local static data with live account progress efficiently, we used **Maps as Lookup Tables**. 

If we have 5,000 achievements in the cache and the account has progress on 200 of them, iterating through the full cache for every account update would be $O(N^2)$. Instead, we load account progress into a map:

```go
progMap := make(map[int]gw2api.AccountAchievement)
for _, a := range accountAch {
    progMap[a.ID] = a
}

// Now lookup is O(1)
for _, cachedAch := range cache.Achievements {
    progress := progMap[cachedAch.ID]
    // ... merge data ...
}
```

### 3. Graceful Degradation
What happens if the cache is missing or the API returns an ID we don't have cached? 

We implemented **Fallbacks**. If a name isn't in the cache, we show the ID:
```go
name := cacheMap[id]
if name == "" {
    name = fmt.Sprintf("Unknown (%d)", id)
}
```
This ensures the app doesn't crash and still provides useful (if slightly less readable) information to the user.

### 4. Categorization via Metadata
We learned that the API's `/v2/achievements` endpoint is "flat." To show the categorized summary view (`./gw2cli achievements`), we had to fetch the `/categories` and `/groups` endpoints and then use those ID lists to filter our merged data. This taught us how to build hierarchical views from flat data sources.
