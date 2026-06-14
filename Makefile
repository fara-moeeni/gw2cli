BINARY_NAME=gw2cli
MAN_SOURCE=docs/man/gw2cli.1.md
MAN_OUTPUT=docs/man/gw2cli.1

build:
	go build -o $(BINARY_NAME) ./cmd/gw2cli

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f $(MAN_OUTPUT)

run:
	go run ./cmd/gw2cli

# Generate man page from markdown (requires go-md2man)
man:
	@if command -v go-md2man > /dev/null; then \
		go-md2man -in $(MAN_SOURCE) -out $(MAN_OUTPUT); \
		echo "Generated $(MAN_OUTPUT)"; \
	else \
		echo "Error: go-md2man not found. Install it with: go install github.com/cpuguy83/go-md2man/v2@latest"; \
		exit 1; \
	fi

# Install binary and man page (may require sudo)
install: build man
	cp $(BINARY_NAME) /usr/local/bin/
	mkdir -p /usr/local/share/man/man1
	cp $(MAN_OUTPUT) /usr/local/share/man/man1/
	@echo "Installed $(BINARY_NAME) and its man page."

.PHONY: build test clean run man install
