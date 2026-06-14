GOFLAGS ?=
LEVELS  ?= ../levels.yml
BIN     := cardinal_chains

.PHONY: all build run test vet fmt clean help

all: build

# Standard optimized build.
build:
	go build $(GOFLAGS) -o $(BIN) .

# Run the game against the bundled levels.yml.
run: build
	./$(BIN) $(LEVELS)

# Run all tests.
test:
	go test ./... -count=1 -v

# Run go vet.
vet:
	go vet ./...

# Format all source files.
fmt:
	go fmt ./...

# Remove build artifacts.
clean:
	rm -f $(BIN)

help:
	@echo "Cardinal Chains (Go TUI) - available targets:"
	@echo "  make build   Compile the binary (default)"
	@echo "  make run     Build, then run the game with ../levels.yml"
	@echo "  make test    Run all tests"
	@echo "  make vet     Run go vet"
	@echo "  make fmt     Format all source files"
	@echo "  make clean   Remove the compiled binary"
