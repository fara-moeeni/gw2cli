# GW2CLI Learning Book

Welcome to the **GW2CLI Learning Book**. Since this project was built via iterative development, this book serves as a historical record and a structured Go learning resource. It distills the lessons learned from building a real-world, production-ready CLI tool from scratch.

As the project evolves, new chapters will be added here to document architectural decisions, new Go concepts learned, and refactoring strategies.

## Chapters

1. [Basics & API Clients](01_basics_and_api.md) - Go modules, Structs, JSON, and HTTP requests.
2. [Concurrency](02_concurrency.md) - Goroutines, WaitGroups, and speeding up API calls.
3. [CLI Flags & Formatting](03_cli_flags_and_formatting.md) - The standard `flag` package and console output.
4. [File I/O & Caching](04_file_io_and_caching.md) - Reading/writing files, creating directories, and cache invalidation.
5. [Architecture & Refactoring](05_architecture_and_refactoring.md) - Moving from flat flags to Subcommands using `flag.FlagSet`.

## How to use this book
If you are reading the codebase and want to understand *why* certain patterns are used or *how* they evolved, refer to the corresponding chapter.
