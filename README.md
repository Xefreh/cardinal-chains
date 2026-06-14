# Cardinal Chains

Cardinal Chains is a puzzle game that challenges players to create sequences by connecting numbers while following specific rules.

## Prerequisites

Before you begin, ensure you have met the following requirements:
- You have a Linux or macOS machine. (Windows steps may vary)
- You have a C compiler installed, such as GCC.
- You have `make` installed.
- You have installed the LibYAML library.

## Installing Cardinal Chains

To install Cardinal Chains, follow these steps:

1. Clone the repository:
   ```
   git clone https://github.com/Xefreh/cardinal-chains.git
   ```
2. Navigate to the project directory:
   ```
   cd cardinal-chains
   ```

## Installing Dependencies

LibYAML is the only external dependency. The bundled script auto-detects your
package manager (`apt-get`, `dnf`, `yum`, `pacman`, `zypper`, `apk`, `brew`)
and installs the right package:

```
make install-deps
# or directly:
./scripts/install-deps.sh
```

To install it manually on Debian/Ubuntu:

```
sudo apt-get install -y libyaml-dev
```

## Compiling the Project

Cardinal Chains uses a Makefile for its build system. To compile the project,
run from the project root:

```
make
```

Other build variants:

```
make debug      # -g -O0 -Wall -Wextra -DDEBUG
make release    # -O3 -DNDEBUG
make help       # list all available targets
```

The `Makefile` honors the standard `CC`, `CFLAGS`, and `LDFLAGS` overrides,
e.g. `make CC=clang CFLAGS="-Wall -Wextra"`.

## Running the Program

The binary expects the path to the YAML levels file as its first argument. The
easiest way to build and run in one step is:

```
make run
```

You can also run the compiled binary directly:

```
./cardinal_chains levels.yml
```

If no argument is given the program prints `Usage: <binary> <YAML file>` and
exits with status 1.

## Go TUI Version

A full terminal user interface (TUI) remake of the game is available in the
[`go/`](go/) directory. It is a 1:1 clone of all game mechanics, built with
[Go](https://go.dev/), [tview](https://github.com/rivo/tview), and
[yaml.v3](https://github.com/go-yaml/yaml). The grid renders in place (no
scrolling), chains are centered and color-coded, and an interactive error bar
notifies you of invalid key presses.

### Prerequisites

- [Go](https://go.dev/dl/) 1.21+ installed.

### Building & Running

From the `go/` directory:

```
cd go
make run        # builds and runs against ../levels.yml
```

Or manually:

```
go build -o cardinal_chains .
./cardinal_chains ../levels.yml
```

### Tests

The Go version includes comprehensive unit tests for every game logic method
(movement rules, cancel/erase/restart, completion check, level advancement)
as well as the level loader, input parser, and renderer:

```
make test       # from the go/ directory
```
