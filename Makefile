CC      ?= gcc
CFLAGS  ?= -Wall -Wextra
LDFLAGS ?=
LDLIBS  ?= -lyaml

SRCS    = $(wildcard src/*.c)
BIN     = cardinal_chains
LEVELS  = levels.yml

.PHONY: all build debug release run clean install-deps help

all: build

# Standard optimized build.
build: $(SRCS)
	$(CC) -O2 $(CFLAGS) $(SRCS) -o $(BIN) $(LDFLAGS) $(LDLIBS)

# Debug build with symbols, kept separate from the default binary.
debug: CFLAGS += -g -O0 -DDEBUG
debug: $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o $(BIN) $(LDFLAGS) $(LDLIBS)

# Release build with full optimization.
release: CFLAGS += -O3 -DNDEBUG
release: $(SRCS)
	$(CC) $(CFLAGS) $(SRCS) -o $(BIN) $(LDFLAGS) $(LDLIBS)

# Build (if needed) and run the game against levels.yml.
run: build
	./$(BIN) $(LEVELS)

# Install build/runtime dependencies (libyaml, etc.).
install-deps:
	@./scripts/install-deps.sh

# Remove build artifacts.
clean:
	rm -f $(BIN)

help:
	@echo "Cardinal Chains - available targets:"
	@echo "  make build        Compile an optimized binary (default)"
	@echo "  make debug        Compile with -g -O0 and DEBUG defined"
	@echo "  make release      Compile with -O3 and NDEBUG defined"
	@echo "  make run          Build, then run the game with levels.yml"
	@echo "  make install-deps Install libyaml and other dependencies"
	@echo "  make clean        Remove the compiled binary"
	@echo "  make help         Show this message"
