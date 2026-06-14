package game

import (
	"cardinal-chains/level"
	"testing"
)

// makeLevel creates a Level from a 2D grid for testing convenience.
func makeLevel(id int, grid [][]int) *level.Level {
	return &level.Level{ID: id, Values: grid}
}

// standardTestGrid is the level from levels.yml:
//
//	-1 1 1 1
//	 0 0 0 0
//	 1 1 1 -1
var standardTestGrid = [][]int{
	{-1, 1, 1, 1},
	{0, 0, 0, 0},
	{1, 1, 1, -1},
}

func newTestGame() *Game {
	lvl := makeLevel(1, standardTestGrid)
	return NewGame(lvl, 0, 3, 4)
}

// ---------------------------------------------------------------------------
// NewGame / initChains
// ---------------------------------------------------------------------------

func TestNewGameChainCount(t *testing.T) {
	g := newTestGame()
	if g.ChainCount() != 2 {
		t.Errorf("ChainCount() = %d, want 2", g.ChainCount())
	}
}

func TestNewGameChainsStartAtAnchors(t *testing.T) {
	g := newTestGame()

	// Chain 0 should start at anchor (0,0) — first -1 found in row-major order.
	if g.Chains[0][0] != (Position{X: 0, Y: 0}) {
		t.Errorf("Chains[0][0] = %v, want {0,0}", g.Chains[0][0])
	}
	// Chain 1 should start at anchor (2,3) — second -1.
	if g.Chains[1][0] != (Position{X: 2, Y: 3}) {
		t.Errorf("Chains[1][0] = %v, want {2,3}", g.Chains[1][0])
	}
}

func TestNewGameInitialChainLengths(t *testing.T) {
	g := newTestGame()
	for i, ch := range g.Chains {
		if len(ch) != 1 {
			t.Errorf("Chains[%d] len = %d, want 1 (anchor only)", i, len(ch))
		}
	}
}

func TestNewGameIsGameOver(t *testing.T) {
	g := newTestGame()
	if g.IsGameOver {
		t.Error("IsGameOver = true, want false for new game")
	}
}

func TestNewGameGridDimensions(t *testing.T) {
	g := newTestGame()
	if g.GridRows() != 3 {
		t.Errorf("GridRows() = %d, want 3", g.GridRows())
	}
	if g.GridCols() != 4 {
		t.Errorf("GridCols() = %d, want 4", g.GridCols())
	}
}

func TestNewGameNoAnchors(t *testing.T) {
	lvl := makeLevel(1, [][]int{{1, 2}, {3, 4}})
	g := NewGame(lvl, 0, 2, 2)
	if g.ChainCount() != 0 {
		t.Errorf("ChainCount() = %d, want 0", g.ChainCount())
	}
}

func TestNewGameMultipleAnchors(t *testing.T) {
	grid := [][]int{
		{-1, 1, -1},
		{1, 0, 1},
		{-1, 1, -1},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 3, 3)

	if g.ChainCount() != 4 {
		t.Fatalf("ChainCount() = %d, want 4", g.ChainCount())
	}

	expected := []Position{{0, 0}, {0, 2}, {2, 0}, {2, 2}}
	for i, want := range expected {
		if g.Chains[i][0] != want {
			t.Errorf("Chains[%d][0] = %v, want %v", i, g.Chains[i][0], want)
		}
	}
}

// ---------------------------------------------------------------------------
// MoveChain — valid moves
// ---------------------------------------------------------------------------

func TestMoveChainEastValid(t *testing.T) {
	g := newTestGame()
	// Chain 0 at (0,0) value -1, move East to (0,1) value 1.
	if !g.MoveChain(0, East) {
		t.Error("MoveChain(0, East) = false, want true")
	}
	if g.Chains[0][1] != (Position{X: 0, Y: 1}) {
		t.Errorf("Chains[0][1] = %v, want {0,1}", g.Chains[0][1])
	}
}

func TestMoveChainWestValid(t *testing.T) {
	g := newTestGame()
	// Chain 1 at (2,3) value -1, move West to (2,2) value 1.
	if !g.MoveChain(1, West) {
		t.Error("MoveChain(1, West) = false, want true")
	}
	if g.Chains[1][1] != (Position{X: 2, Y: 2}) {
		t.Errorf("Chains[1][1] = %v, want {2,2}", g.Chains[1][1])
	}
}

func TestMoveChainSouthValid(t *testing.T) {
	grid := [][]int{
		{-1, 0, 0},
		{1, 0, 0},
		{2, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 3, 3)

	if !g.MoveChain(0, South) {
		t.Error("MoveChain(0, South) = false, want true")
	}
	if g.Chains[0][1] != (Position{X: 1, Y: 0}) {
		t.Errorf("Chains[0][1] = %v, want {1,0}", g.Chains[0][1])
	}
}

func TestMoveChainNorthValid(t *testing.T) {
	grid := [][]int{
		{3, 0, 0},
		{2, 0, 0},
		{-1, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 3, 3)

	// From (2,0) move North to (1,0) value 2.
	if !g.MoveChain(0, North) {
		t.Error("MoveChain(0, North) = false, want true")
	}
	if g.Chains[0][1] != (Position{X: 1, Y: 0}) {
		t.Errorf("Chains[0][1] = %v, want {1,0}", g.Chains[0][1])
	}
}

func TestMoveChainAnchorToPositiveValue(t *testing.T) {
	g := newTestGame()
	// Anchor value is -1. Moving to a cell with value 1 is legal (1 >= -1).
	if !g.MoveChain(0, East) {
		t.Error("Move from anchor (-1) to value 1 should succeed")
	}
}

func TestMoveChainToEqualValue(t *testing.T) {
	g := newTestGame()
	// Move chain 0: (0,0) → (0,1) → (0,2) — all values 1 (non-decreasing).
	g.MoveChain(0, East) // -1 → 1
	if !g.MoveChain(0, East) {
		t.Error("Move from value 1 to value 1 should succeed (non-decreasing)")
	}
}

func TestMoveChainToHigherValue(t *testing.T) {
	grid := [][]int{
		{-1, 2, 3},
		{0, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 2, 3)

	// -1 → 2 (legal)
	if !g.MoveChain(0, East) {
		t.Error("Move from -1 to 2 should succeed")
	}
	// 2 → 3 (legal, non-decreasing)
	if !g.MoveChain(0, East) {
		t.Error("Move from 2 to 3 should succeed")
	}
}

// ---------------------------------------------------------------------------
// MoveChain — out of bounds
// ---------------------------------------------------------------------------

func TestMoveChainOutOfBoundsNorth(t *testing.T) {
	g := newTestGame()
	// Chain 0 at row 0, moving North goes out of bounds.
	if g.MoveChain(0, North) {
		t.Error("MoveChain(0, North) should fail (out of bounds at row 0)")
	}
}

func TestMoveChainOutOfBoundsWest(t *testing.T) {
	g := newTestGame()
	// Chain 0 at col 0, moving West goes out of bounds.
	if g.MoveChain(0, West) {
		t.Error("MoveChain(0, West) should fail (out of bounds at col 0)")
	}
}

func TestMoveChainOutOfBoundsSouth(t *testing.T) {
	g := newTestGame()
	// Chain 1 at row 2 (last row), moving South goes out of bounds.
	if g.MoveChain(1, South) {
		t.Error("MoveChain(1, South) should fail (out of bounds at last row)")
	}
}

func TestMoveChainOutOfBoundsEast(t *testing.T) {
	g := newTestGame()
	// Chain 1 at col 3 (last col), moving East goes out of bounds.
	if g.MoveChain(1, East) {
		t.Error("MoveChain(1, East) should fail (out of bounds at last col)")
	}
}

// ---------------------------------------------------------------------------
// MoveChain — onto occupied cell
// ---------------------------------------------------------------------------

func TestMoveChainOntoOccupiedCell(t *testing.T) {
	// Build a grid where two chains can potentially reach the same cell.
	grid := [][]int{
		{-1, 1, -1},
		{0, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 2, 3)

	// Chain 0 at (0,0), chain 1 at (0,2).
	// Move chain 0 East to (0,1).
	g.MoveChain(0, East)
	// Now try to move chain 1 West to (0,1) — it's occupied.
	if g.MoveChain(1, West) {
		t.Error("MoveChain(1, West) should fail — (0,1) is occupied by chain 0")
	}
}

func TestMoveChainOntoOwnCell(t *testing.T) {
	g := newTestGame()
	// Chain 0 at (0,0). Move East to (0,1).
	g.MoveChain(0, East)
	// Try to move West back to (0,0) — it's occupied by the chain's own anchor.
	if g.MoveChain(0, West) {
		t.Error("MoveChain(0, West) should fail — (0,0) is occupied by own anchor")
	}
}

// ---------------------------------------------------------------------------
// MoveChain — value rules
// ---------------------------------------------------------------------------

func TestMoveChainOntoZeroCell(t *testing.T) {
	g := newTestGame()
	// Chain 0 at (0,0). Moving South goes to (1,0) which has value 0.
	if g.MoveChain(0, South) {
		t.Error("MoveChain(0, South) should fail — target cell value is 0")
	}
}

func TestMoveChainDecreasingValue(t *testing.T) {
	grid := [][]int{
		{3, 2, 1},
		{-1, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 2, 3)

	// Anchor at (1,0) value -1. Move North to (0,0) value 3: legal (3 >= -1).
	g.MoveChain(0, North)
	// Now at (0,0) value 3. Move East to (0,1) value 2: illegal (2 < 3).
	if g.MoveChain(0, East) {
		t.Error("Move from value 3 to value 2 should fail (decreasing)")
	}
}

func TestMoveChainFromZeroAnchorBlockedByZero(t *testing.T) {
	grid := [][]int{
		{-1, 0, 1},
		{0, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 2, 3)

	// Anchor at (0,0) value -1. Move East to (0,1) value 0: blocked.
	if g.MoveChain(0, East) {
		t.Error("Move onto zero cell should fail")
	}
}

// ---------------------------------------------------------------------------
// MoveChain — invalid chain index
// ---------------------------------------------------------------------------

func TestMoveChainInvalidIndexNegative(t *testing.T) {
	g := newTestGame()
	if g.MoveChain(-1, East) {
		t.Error("MoveChain(-1, ...) should fail (invalid index)")
	}
}

func TestMoveChainInvalidIndexTooLarge(t *testing.T) {
	g := newTestGame()
	if g.MoveChain(99, East) {
		t.Error("MoveChain(99, ...) should fail (invalid index)")
	}
}

// ---------------------------------------------------------------------------
// MoveChain — chain length updates
// ---------------------------------------------------------------------------

func TestMoveChainIncrementsLength(t *testing.T) {
	g := newTestGame()
	if g.ChainLength(0) != 1 {
		t.Fatalf("Initial ChainLength(0) = %d, want 1", g.ChainLength(0))
	}
	g.MoveChain(0, East)
	if g.ChainLength(0) != 2 {
		t.Errorf("After move, ChainLength(0) = %d, want 2", g.ChainLength(0))
	}
	g.MoveChain(0, East)
	if g.ChainLength(0) != 3 {
		t.Errorf("After 2 moves, ChainLength(0) = %d, want 3", g.ChainLength(0))
	}
}

func TestMoveChainFailedMoveDoesNotChangeLength(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, North) // fails (out of bounds)
	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1 (failed move should not extend)", g.ChainLength(0))
	}
}

// ---------------------------------------------------------------------------
// CancelLastMove
// ---------------------------------------------------------------------------

func TestCancelLastMove(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	if g.ChainLength(0) != 3 {
		t.Fatalf("ChainLength(0) = %d, want 3 before cancel", g.ChainLength(0))
	}

	g.CancelLastMove(0)
	if g.ChainLength(0) != 2 {
		t.Errorf("ChainLength(0) = %d, want 2 after cancel", g.ChainLength(0))
	}

	last := g.Chains[0][g.ChainLength(0)-1]
	if last != (Position{X: 0, Y: 1}) {
		t.Errorf("Last position = %v, want {0,1}", last)
	}
}

func TestCancelLastMoveKeepsAnchor(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)

	g.CancelLastMove(0) // back to anchor
	if g.ChainLength(0) != 1 {
		t.Fatalf("ChainLength(0) = %d, want 1 after cancel", g.ChainLength(0))
	}

	// Cancel again — should be a no-op (chain can't go below 1).
	g.CancelLastMove(0)
	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1 (cancel on anchor-only chain is no-op)", g.ChainLength(0))
	}
}

func TestCancelLastMoveOnAnchorOnly(t *testing.T) {
	g := newTestGame()
	g.CancelLastMove(0)
	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1", g.ChainLength(0))
	}
}

func TestCancelLastMoveInvalidIndex(t *testing.T) {
	g := newTestGame()
	g.CancelLastMove(-1) // should not panic
	g.CancelLastMove(99) // should not panic
}

// ---------------------------------------------------------------------------
// EraseChain
// ---------------------------------------------------------------------------

func TestEraseChain(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	if g.ChainLength(0) != 4 {
		t.Fatalf("ChainLength(0) = %d, want 4 before erase", g.ChainLength(0))
	}

	g.EraseChain(0)
	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1 after erase", g.ChainLength(0))
	}
	if g.Chains[0][0] != (Position{X: 0, Y: 0}) {
		t.Errorf("After erase, anchor = %v, want {0,0}", g.Chains[0][0])
	}
}

func TestEraseChainOnAnchorOnly(t *testing.T) {
	g := newTestGame()
	g.EraseChain(0)
	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1 (erase on anchor-only chain is no-op)", g.ChainLength(0))
	}
}

func TestEraseChainInvalidIndex(t *testing.T) {
	g := newTestGame()
	g.EraseChain(-1) // should not panic
	g.EraseChain(99) // should not panic
}

func TestEraseChainOnlyAffectsTargetChain(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)
	g.MoveChain(1, West)

	g.EraseChain(0)

	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1", g.ChainLength(0))
	}
	if g.ChainLength(1) != 2 {
		t.Errorf("ChainLength(1) = %d, want 2 (unaffected by erase of chain 0)", g.ChainLength(1))
	}
}

// ---------------------------------------------------------------------------
// RestartLevel
// ---------------------------------------------------------------------------

func TestRestartLevel(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	g.MoveChain(1, West)

	g.RestartLevel()

	for i := range g.Chains {
		if len(g.Chains[i]) != 1 {
			t.Errorf("Chain %d length = %d, want 1 after restart", i, len(g.Chains[i]))
		}
	}
}

func TestRestartLevelPreservesAnchors(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East)
	g.MoveChain(1, West)

	g.RestartLevel()

	if g.Chains[0][0] != (Position{X: 0, Y: 0}) {
		t.Errorf("Chain 0 anchor = %v, want {0,0}", g.Chains[0][0])
	}
	if g.Chains[1][0] != (Position{X: 2, Y: 3}) {
		t.Errorf("Chain 1 anchor = %v, want {2,3}", g.Chains[1][0])
	}
}

func TestRestartLevelOnFreshGame(t *testing.T) {
	g := newTestGame()
	g.RestartLevel()

	for i := range g.Chains {
		if len(g.Chains[i]) != 1 {
			t.Errorf("Chain %d length = %d, want 1", i, len(g.Chains[i]))
		}
	}
}

// ---------------------------------------------------------------------------
// IsGameCompleted
// ---------------------------------------------------------------------------

func TestIsGameCompletedInitiallyFalse(t *testing.T) {
	g := newTestGame()
	if g.IsGameCompleted() {
		t.Error("IsGameCompleted() = true, want false for new game")
	}
}

func TestIsGameCompletedPartiallyFilled(t *testing.T) {
	g := newTestGame()
	g.MoveChain(0, East) // only (0,1) is covered by chain 0
	if g.IsGameCompleted() {
		t.Error("IsGameCompleted() = true, want false (partially filled)")
	}
}

func TestIsGameCompletedTrue(t *testing.T) {
	g := newTestGame()

	// Chain 0: (0,0) → (0,1) → (0,2) → (0,3)
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	g.MoveChain(0, East)

	// Chain 1: (2,3) → (2,2) → (2,1) → (2,0)
	g.MoveChain(1, West)
	g.MoveChain(1, West)
	g.MoveChain(1, West)

	if !g.IsGameCompleted() {
		t.Error("IsGameCompleted() = false, want true (all positive cells covered)")
	}
}

func TestIsGameCompletedNoPositiveCells(t *testing.T) {
	// A grid with only anchors and zeros has no positive cells to cover.
	grid := [][]int{
		{-1, 0, -1},
		{0, 0, 0},
	}
	lvl := makeLevel(1, grid)
	g := NewGame(lvl, 0, 2, 3)

	if !g.IsGameCompleted() {
		t.Error("IsGameCompleted() = false, want true (no positive cells to cover)")
	}
}

func TestIsGameCompletedOneChainMissing(t *testing.T) {
	g := newTestGame()

	// Complete chain 0 fully, but leave chain 1 at anchor.
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	g.MoveChain(0, East)

	if g.IsGameCompleted() {
		t.Error("IsGameCompleted() = true, want false (chain 1 not extended)")
	}
}

// ---------------------------------------------------------------------------
// LoadNextLevel
// ---------------------------------------------------------------------------

func TestLoadNextLevel(t *testing.T) {
	grid1 := [][]int{{-1, 1}, {0, 0}}
	grid2 := [][]int{{-1, 2, 2}, {0, 0, 0}}

	ls := level.NewLevels()
	ls.AddLevel(1, grid1)
	ls.AddLevel(2, grid2)

	g := NewGame(&ls.Items[0], 0, 2, 3)
	g.CurrentLevel = 1

	g.LoadNextLevel(ls)

	if g.Level.ID != 2 {
		t.Errorf("Level.ID = %d, want 2", g.Level.ID)
	}
	if g.ChainCount() != 1 {
		t.Errorf("ChainCount() = %d, want 1", g.ChainCount())
	}
	if g.Chains[0][0] != (Position{X: 0, Y: 0}) {
		t.Errorf("Chains[0][0] = %v, want {0,0}", g.Chains[0][0])
	}
	if g.IsGameOver {
		t.Error("IsGameOver = true, want false")
	}
}

func TestLoadNextLevelSetsGameOver(t *testing.T) {
	grid1 := [][]int{{-1, 1}, {0, 0}}

	ls := level.NewLevels()
	ls.AddLevel(1, grid1)

	g := NewGame(&ls.Items[0], 0, 2, 2)
	g.CurrentLevel = 1 // beyond levels count

	g.LoadNextLevel(ls)

	if !g.IsGameOver {
		t.Error("IsGameOver = false, want true (no more levels)")
	}
}

func TestLoadNextLevelResetsChains(t *testing.T) {
	grid1 := [][]int{{-1, 1, 1}, {0, 0, 0}}
	grid2 := [][]int{{-1, 1, 1}, {0, 0, 0}}

	ls := level.NewLevels()
	ls.AddLevel(1, grid1)
	ls.AddLevel(2, grid2)

	g := NewGame(&ls.Items[0], 0, 2, 3)
	g.MoveChain(0, East)
	g.MoveChain(0, East)
	if g.ChainLength(0) != 3 {
		t.Fatalf("ChainLength(0) = %d, want 3 before load", g.ChainLength(0))
	}

	g.CurrentLevel = 1
	g.LoadNextLevel(ls)

	if g.ChainLength(0) != 1 {
		t.Errorf("ChainLength(0) = %d, want 1 after load (chains reset)", g.ChainLength(0))
	}
}

// ---------------------------------------------------------------------------
// ChainLength / ChainCount
// ---------------------------------------------------------------------------

func TestChainLengthInvalidIndex(t *testing.T) {
	g := newTestGame()
	if g.ChainLength(-1) != 0 {
		t.Errorf("ChainLength(-1) = %d, want 0", g.ChainLength(-1))
	}
	if g.ChainLength(99) != 0 {
		t.Errorf("ChainLength(99) = %d, want 0", g.ChainLength(99))
	}
}

// ---------------------------------------------------------------------------
// Full game scenario — solving the standard level
// ---------------------------------------------------------------------------

func TestFullGameScenario(t *testing.T) {
	g := newTestGame()

	// Step 1: Extend chain 0 eastward to cover the top row.
	if !g.MoveChain(0, East) {
		t.Fatal("Move 1 failed")
	}
	if !g.MoveChain(0, East) {
		t.Fatal("Move 2 failed")
	}
	if !g.MoveChain(0, East) {
		t.Fatal("Move 3 failed")
	}

	// Not yet completed — chain 1 hasn't moved.
	if g.IsGameCompleted() {
		t.Fatal("Game should not be completed yet")
	}

	// Step 2: Extend chain 1 westward to cover the bottom row.
	if !g.MoveChain(1, West) {
		t.Fatal("Move 4 failed")
	}
	if !g.MoveChain(1, West) {
		t.Fatal("Move 5 failed")
	}
	if !g.MoveChain(1, West) {
		t.Fatal("Move 6 failed")
	}

	// Now all positive cells are covered.
	if !g.IsGameCompleted() {
		t.Fatal("Game should be completed")
	}

	// Step 3: Undo and verify not completed.
	g.CancelLastMove(1)
	if g.IsGameCompleted() {
		t.Fatal("Game should not be completed after undo")
	}

	// Step 4: Restart and verify not completed.
	g.RestartLevel()
	if g.IsGameCompleted() {
		t.Fatal("Game should not be completed after restart")
	}
	for i := range g.Chains {
		if len(g.Chains[i]) != 1 {
			t.Errorf("Chain %d length = %d, want 1 after restart", i, len(g.Chains[i]))
		}
	}
}
