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
- `CMakeLists.txt` — alternative CMake build. Note that it defines an
  `INPUT_FILE` macro and copies `levels.yml` into the build dir, but the
  current `main.c` reads the path from `argv[1]` instead, so that macro is
  effectively unused.

## Build with gcc

The only external dependency is LibYAML (`libyaml-dev` / `yaml.h`). On Debian/
Ubuntu install it with:

```
sudo apt-get install -y libyaml-dev
```

Compile directly with gcc (this is the recommended workflow for this project):

```
gcc main.c -o cardinal_chains -lyaml
```

For a debug build with warnings and symbols:

```
gcc -Wall -Wextra -g main.c -o cardinal_chains -lyaml
```

For an optimized release build:

```
gcc -O2 main.c -o cardinal_chains -lyaml
```

## Run

The binary expects the path to the YAML levels file as its first argument:

```
./cardinal_chains levels.yml
```

If no argument is given the program prints `Usage: <binary> <YAML file>` and
exits with status 1.

## Build with CMake (alternative)

```
mkdir build && cd build
cmake ..
make
./cardinal_chains ../levels.yml
```

Note: the CMake target also accepts the YAML path as `argv[1]`; run it from the
`build/` directory and point it at `../levels.yml` (or the copied
`levels.yml`).

## Verification

There is no test suite in the repository. After building, verify the binary
works by running it against the bundled level file and confirming the grid is
rendered and input is accepted:

```
./cardinal_chains levels.yml
```
