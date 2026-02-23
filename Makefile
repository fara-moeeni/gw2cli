BINARY_NAME=gw2cli

build:
	go build -o $(BINARY_NAME) cmd/gw2cli/main.go

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

run:
	go run cmd/gw2cli/main.go

.PHONY: build test clean run
