BINARY=ollama-mcp
GO=go
LDFLAGS=-ldflags="-w -s"
BUILD_ENV=CGO_ENABLED=0

.PHONY: all build clean lint lint-docs test fmt

all: build

build: $(BINARY)

$(BINARY): $(wildcard *.go)
	$(BUILD_ENV) $(GO) build $(LDFLAGS) -o $(BINARY)

clean:
	rm -f $(BINARY)

fmt:
	$(GO) fmt ./...

lint:
	$(GO) vet ./...
	@if command -v staticcheck >/dev/null; then staticcheck ./...; else echo "staticcheck not found, skipping..."; fi

lint-docs:
	@echo "Linting manpage..."
	@if command -v mandoc >/dev/null; then mandoc -Tlint ollama-mcp.1; else echo "mandoc not found, skipping..."; fi
	@echo "Linting markdown documentation..."
	@if command -v markdownlint-cli2 >/dev/null; then markdownlint-cli2 README.md DESIGN.md; else echo "markdownlint-cli2 not found, skipping..."; fi

test:
	$(GO) test ./...
