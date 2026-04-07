# GW2CLI Learning Book

Welcome to the **GW2CLI Learning Book**. Since this project was built via iterative development, this book serves as a historical record and a structured Go learning resource. It distills the lessons learned from building a real-world, production-ready CLI tool from scratch.

As the project evolves, new chapters will be added here to document architectural decisions, new Go concepts learned, and refactoring strategies.

## Chapters

### The Go Crash Course (For Beginners)
0. [Go Crash Course](00_go_crash_course.md) - Basic syntax, variables, collections, and structs.
0.1. [Pointers & References](00_1_pointers_and_references.md) - Memory addresses, passing by reference, and why the `flag` package uses them.
0.2. [Interfaces & Errors](00_2_interfaces_and_errors.md) - Understanding interfaces, polymorphism, and how Go handles errors.

### Building GW2CLI
1. [Basics & API Clients](01_basics_and_api.md) - Go modules, Structs, JSON, and HTTP requests.
2. [Concurrency](02_concurrency.md) - Goroutines, WaitGroups, and speeding up API calls.
3. [CLI Flags & Formatting](03_cli_flags_and_formatting.md) - The standard `flag` package and console output.
4. [File I/O & Caching](04_file_io_and_caching.md) - Reading/writing files, creating directories, and cache invalidation.
5. [Architecture & Refactoring](05_architecture_and_refactoring.md) - Moving from flat flags to Subcommands using `flag.FlagSet`.
6. [Advanced Batching & The Interface Trap](06_advanced_batching_and_generics.md) - Why generic ID resolution failed and how we fixed it.
7. [Tracking the Tick-Tock (Resets & States)](07_time_and_reset_logic.md) - Time handling, resets, and completion status.
8. [Data Merging & Multi-Caching](08_data_merging_and_multi_caching.md) - Combining static local data with live API state.

## How to use this book
If you are reading the codebase and want to understand *why* certain patterns are used or *how* they evolved, refer to the corresponding chapter.
