# Chapter 7: Tracking the Tick-Tock (Resets & States)

## Context
Phase 18 introduced the "Daily & Weekly Tracker." This required us to handle two types of temporal data: 
1. **Public Dailies**: What achievements are available for everyone today?
2. **Account Progress**: Which of those (and other bosses/dungeons) has the user actually finished?

The challenge was presenting a unified view of "State" (Done vs. Not Done) and clearly communicating when that state would expire (the Reset).

## Key Go Concepts Learned

### 1. Unified Status Formatting
To make the CLI feel like a professional tool, we needed a consistent way to show "Done" across different subcommands. We adopted the `[✓]` and `[ ]` pattern.

In `internal/ui/printer.go`, we used simple conditional logic to handle this:
```go
status := "[ ]"
if item.Completed {
    status = "[✓]"
}
fmt.Printf(" %s %s\n", status, item.Name)
```

### 2. Aggregating Data from Multiple Sources
The "Daily" view is actually a composite of 4 different API calls:
- World Bosses (Account)
- Dungeons (Account)
- Fractal Achievements (Public)
- Wizard's Vault (Account)

We learned to use the `inventory.Service` to fetch these in parallel or sequence, and then pass them into specialized "Print" functions. This kept the `main.go` routing logic clean while allowing the UI layer to focus purely on layout.

### 3. Handling API Limitations (The Convergences Case)
A critical lesson in real-world API consumption: **The API doesn't always have what you need.** 

The weekly Convergences boss rotation is not exposed by the GW2 API. Instead of trying to "guess" the rotation or hardcode a calendar (which would break eventually), we chose **Transparency**. 

We added a hardcoded note in the `weekly` output:
```text
Convergences: rotation not available via API — check https://wiki.guildwars2.com/wiki/Convergences
```
*Lesson Learned: If you can't provide data accurately, it's better to admit it and point the user to a reliable source than to provide potentially wrong information.*

### 4. Reset Timers
We learned that clearly stating the reset time (`Daily reset: 00:00 UTC`) is essential context for the user. Without it, a `[✓]` status might be confusing if the user doesn't know when it's going to flip back to `[ ]`.
