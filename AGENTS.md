# AGENTS.md

## Project Goal

Cardinal Chains is a terminal-based puzzle game written in C. The player must
fill a 2D number grid by drawing "chains" that connect anchor cells, following
strict movement rules. Each level is defined in `levels.yml` and parsed at
runtime with LibYAML.

### How a level works
- The grid is read from the YAML file passed as the first command-line argument.
- Cell values have special meaning:
  - `-1` — an anchor cell. Anchors always come in pairs and form the start/end
    of a chain. Each `-1` becomes the seed of one chain.
  - `0` — an empty/blocked cell that cannot be entered by any chain.
  - Positive integers — fillable cells. Every positive cell must be covered by
    some chain for the level to be considered complete.
- The player picks a chain and moves it one step at a time in a cardinal
  direction (North/South/East/West). A move is only legal if:
  1. It stays inside the grid bounds.
  2. The target cell is not already occupied by any chain.
  3. The target cell is non-zero and its value is greater than or equal to the
     value of the cell the chain is currently on (chains are non-decreasing in
     value as they grow).
- A level is solved when every positive-value cell in the grid is covered by a
  chain. The game then advances to the next level.

### Interactive commands (entered at the prompt)
| Key | Action |
| --- | ------ |
| `N` / `S` / `E` / `W` | Move the current chain north/south/east/west |
| `B` | Cancel (undo) the previous move on the current chain |
| `R` | Erase the current chain back to its anchor |
| `X` | Restart the whole level (erase every chain) |
| `C` | Cycle selection to the next chain |
| `Q` | Quit the game immediately |

Chains are rendered with ANSI colors so a color-capable terminal is expected.

## Source layout
- `main.c` — the entire game: YAML parsing, grid model, chain logic, input loop,
  rendering, and memory cleanup.
- `levels.yml` — the level data. Each top-level entry under `levels:` maps a
  level number to a list of rows; each row is a space-separated string of ints.
- `Makefile` — the build system. Drives compilation, running, and cleanup.
- `scripts/install-deps.sh` — detects the system package manager and installs
  libyaml and any other build/runtime dependencies.

## Dependencies

The only external dependency is LibYAML (`libyaml-dev` / `yaml.h`).

## Install dependencies

The install script auto-detects the system package manager (`apt-get`, `dnf`,
`yum`, `pacman`, `zypper`, `apk`, `brew`) and installs the right libyaml
package. It re-execs itself under `sudo` if not root:

```
make install-deps
# or directly:
./scripts/install-deps.sh
```

To install manually on Debian/Ubuntu:

```
sudo apt-get install -y libyaml-dev
```

## Build

The project uses a Makefile as its build system (CMake is no longer used).
`make` / `make build` produces an optimized `cardinal_chains` binary using gcc
and `-lyaml`:

```
make            # == make build
```

Other build targets:

```
make debug      # -g -O0 -Wall -Wextra -DDEBUG
make release    # -O3 -DNDEBUG
make help       # list all available targets
```

The `Makefile` honors the standard `CC`, `CFLAGS`, `LDFLAGS` overrides, e.g.
`make CC=clang CFLAGS="-Wall -Wextra"`.

For a plain gcc invocation without make (equivalent to `make build`):

```
gcc main.c -o cardinal_chains -lyaml
```

## Run

`make run` builds (if needed) and runs the game against `levels.yml`:

```
make run
```

The binary also expects the path to the YAML levels file as its first argument,
so it can be invoked directly:

```
./cardinal_chains levels.yml
```

If no argument is given the program prints `Usage: <binary> <YAML file>` and
exits with status 1.

## Other targets

```
make clean        # remove the compiled binary
make install-deps # install libyaml via scripts/install-deps.sh
```

## Verification

There is no test suite in the repository. After building, verify the binary
works by running it against the bundled level file and confirming the grid is
rendered and input is accepted:

```
make run
# or: ./cardinal_chains levels.yml
```
