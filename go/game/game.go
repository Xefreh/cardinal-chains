package game

import "cardinal-chains/level"

// Position represents a cell coordinate in the grid.
// X is the row index, Y is the column index (matching the C Position struct).
type Position struct {
	X int
	Y int
}

// Direction represents one of the four cardinal directions.
type Direction int

const (
	North Direction = iota
	South
	East
	West
)

// Game holds the full state of a Cardinal Chains game session.
//
// It mirrors the C CardinalChainsGame struct:
//   - Level: the current level being played.
//   - Chains: a slice of chains, each chain is a slice of Positions.
//   - IsGameOver: true when all levels have been completed.
//   - CurrentLevel: the index into the Levels collection for the current level.
type Game struct {
	Level        level.Level
	Chains       [][]Position
	IsGameOver   bool
	CurrentLevel int

	// gridRows and gridCols are the global maximum dimensions across all
	// levels, used for movement bounds-checking (mirroring main.c).
	gridRows int
	gridCols int
}

// NewGame creates a new game initialized with the given level. The
// currentLevel index and global grid dimensions are also stored.
func NewGame(lvl *level.Level, currentLevel, gridRows, gridCols int) *Game {
	g := &Game{
		Level:        *lvl,
		IsGameOver:   false,
		CurrentLevel: currentLevel,
		gridRows:     gridRows,
		gridCols:     gridCols,
	}
	g.initChains()
	return g
}

// findNextEmptyChain returns the index of the first chain whose length is 0,
// or -1 if none is found. Used during initialization to assign anchors to
// chains in order.
func (g *Game) findNextEmptyChain() int {
	for i, ch := range g.Chains {
		if len(ch) == 0 {
			return i
		}
	}
	return -1
}

// initChains initializes the chains slice based on the level's anchor cells.
// Each anchor cell (value == -1) becomes the seed (first Position) of one
// chain. Anchors are found in row-major order, matching the C implementation.
func (g *Game) initChains() {
	chainCount := g.Level.CountChains()
	g.Chains = make([][]Position, chainCount)

	for i := range g.Level.Values {
		for j := range g.Level.Values[i] {
			if g.Level.Values[i][j] == -1 {
				idx := g.findNextEmptyChain()
				if idx == -1 {
					return
				}
				g.Chains[idx] = []Position{{X: i, Y: j}}
			}
		}
	}
}

// GridRows returns the global grid row count used for bounds checking.
func (g *Game) GridRows() int {
	return g.gridRows
}

// GridCols returns the global grid column count used for bounds checking.
func (g *Game) GridCols() int {
	return g.gridCols
}

// ChainCount returns the number of chains in the game.
func (g *Game) ChainCount() int {
	return len(g.Chains)
}

// ChainLength returns the length (number of positions) of the given chain.
func (g *Game) ChainLength(chainIndex int) int {
	if chainIndex < 0 || chainIndex >= len(g.Chains) {
		return 0
	}
	return len(g.Chains[chainIndex])
}

// MoveChain attempts to extend the given chain by one step in the specified
// direction. Returns true if the move was legal and applied, false otherwise.
//
// A move is legal if and only if:
//  1. The chain index is valid.
//  2. The target cell is inside the grid bounds.
//  3. The target cell is not already occupied by any chain.
//  4. The target cell's value is non-zero and greater than or equal to the
//     value of the cell the chain is currently on (non-decreasing values).
func (g *Game) MoveChain(chainIndex int, dir Direction) bool {
	if chainIndex < 0 || chainIndex >= len(g.Chains) {
		return false
	}

	chain := g.Chains[chainIndex]
	lastPos := chain[len(chain)-1]
	newPos := lastPos

	switch dir {
	case North:
		if lastPos.X <= 0 {
			return false
		}
		newPos.X--
	case South:
		if lastPos.X >= g.gridRows-1 {
			return false
		}
		newPos.X++
	case East:
		if lastPos.Y >= g.gridCols-1 {
			return false
		}
		newPos.Y++
	case West:
		if lastPos.Y <= 0 {
			return false
		}
		newPos.Y--
	}

	for _, ch := range g.Chains {
		for _, pos := range ch {
			if pos.X == newPos.X && pos.Y == newPos.Y {
				return false
			}
		}
	}

	lastMoveValue := g.Level.CellValue(lastPos.X, lastPos.Y)
	nextMoveValue := g.Level.CellValue(newPos.X, newPos.Y)

	if nextMoveValue < lastMoveValue || nextMoveValue == 0 {
		return false
	}

	g.Chains[chainIndex] = append(chain, newPos)
	return true
}

// CancelLastMove removes the last position from the given chain. The anchor
// (first position) is always kept — the chain length will never drop below 1.
func (g *Game) CancelLastMove(chainIndex int) {
	if chainIndex < 0 || chainIndex >= len(g.Chains) {
		return
	}
	chain := g.Chains[chainIndex]
	if len(chain) > 1 {
		g.Chains[chainIndex] = chain[:len(chain)-1]
	}
}

// EraseChain removes all positions from the given chain except the anchor
// (the first position), resetting it back to length 1.
func (g *Game) EraseChain(chainIndex int) {
	if chainIndex < 0 || chainIndex >= len(g.Chains) {
		return
	}
	if len(g.Chains[chainIndex]) > 1 {
		g.Chains[chainIndex] = g.Chains[chainIndex][:1]
	}
}

// RestartLevel erases every chain back to its anchor.
func (g *Game) RestartLevel() {
	for i := range g.Chains {
		g.EraseChain(i)
	}
}

// isCellFilled checks whether the cell at (row, col) is covered by any chain.
func (g *Game) isCellFilled(row, col int) bool {
	for _, ch := range g.Chains {
		for _, pos := range ch {
			if pos.X == row && pos.Y == col {
				return true
			}
		}
	}
	return false
}

// IsGameCompleted returns true when every positive-value cell in the grid is
// covered by some chain. Anchor cells (-1) and blocked cells (0) are ignored.
func (g *Game) IsGameCompleted() bool {
	for i := range g.Level.Values {
		for j := range g.Level.Values[i] {
			v := g.Level.Values[i][j]
			if v == -1 || v == 0 {
				continue
			}
			if !g.isCellFilled(i, j) {
				return false
			}
		}
	}
	return true
}

// LoadNextLevel advances to the next level in the collection. If there are
// no more levels, IsGameOver is set to true.
func (g *Game) LoadNextLevel(levels *level.Levels) {
	if g.CurrentLevel >= levels.Count() {
		g.IsGameOver = true
		return
	}
	g.Level = levels.Items[g.CurrentLevel]
	g.initChains()
}
