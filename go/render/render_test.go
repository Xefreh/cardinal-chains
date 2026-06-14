package render

import (
	"cardinal-chains/game"
	"cardinal-chains/level"
	"strings"
	"testing"
)

func makeGame() *game.Game {
	lvl := &level.Level{ID: 1, Values: [][]int{
		{-1, 1, 1, 1},
		{0, 0, 0, 0},
		{1, 1, 1, -1},
	}}
	return game.NewGame(lvl, 0, 3, 4)
}

func TestGridStringBasic(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}
	output := GridString(g, 0, activeChains)

	if output == "" {
		t.Error("GridString returned empty string")
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(lines))
	}
}

func TestGridStringShowsAnchorsAsX(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}
	output := GridString(g, 0, activeChains)

	if !strings.Contains(output, "x ") {
		t.Errorf("Expected 'x' for anchor cells in output:\n%s", output)
	}
}

func TestGridStringShowsBlankForZero(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}
	output := GridString(g, 0, activeChains)

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 2 {
		t.Fatalf("Expected at least 2 lines, got %d", len(lines))
	}

	// Row index 1 (the zero row) should contain blank spaces, not "0".
	zeroRow := lines[1]
	if strings.Contains(zeroRow, "0 ") && !strings.Contains(zeroRow, "[") {
		// Check that it's all blanks (no digits). Each zero cell renders as "  ".
		// The entire row should be whitespace.
		trimmed := strings.TrimSpace(zeroRow)
		if trimmed != "" {
			t.Errorf("Expected blank row for zeros, got %q", zeroRow)
		}
	}
}

func TestGridStringContainsColorTags(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}

	// Move a chain so it has colored cells.
	g.MoveChain(0, game.East)

	output := GridString(g, 0, activeChains)
	if !strings.Contains(output, "[red]") {
		t.Errorf("Expected '[red]' color tag in output:\n%s", output)
	}
	if !strings.Contains(output, "[-:]") {
		t.Errorf("Expected color reset tag '[-:]' in output:\n%s", output)
	}
}

func TestStatusString(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}
	output := StatusString(g, 0, activeChains)

	if !strings.Contains(output, "Current position:") {
		t.Errorf("Expected 'Current position:' in status:\n%s", output)
	}
	if !strings.Contains(output, "Current chain color:") {
		t.Errorf("Expected 'Current chain color:' in status:\n%s", output)
	}
}

func TestStatusStringShowsCorrectPosition(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}

	g.MoveChain(0, game.East) // chain 0 now at (0,1)

	output := StatusString(g, 0, activeChains)
	// Position should be row 1 (0+1), column 2 (1+1) in 1-indexed display.
	if !strings.Contains(output, "row 1") {
		t.Errorf("Expected 'row 1' in status:\n%s", output)
	}
	if !strings.Contains(output, "column 2") {
		t.Errorf("Expected 'column 2' in status:\n%s", output)
	}
}

func TestHelpString(t *testing.T) {
	output := HelpString()
	expectedSubstrings := []string{
		"N/S/E/W",
		"B",
		"R",
		"X",
		"C",
		"Q",
	}
	for _, s := range expectedSubstrings {
		if !strings.Contains(output, s) {
			t.Errorf("Expected %q in help string:\n%s", s, output)
		}
	}
}

func TestFullRender(t *testing.T) {
	g := makeGame()
	activeChains := []int{0, 1}
	output := FullRender(g, 0, activeChains)

	if !strings.Contains(output, "Current position:") {
		t.Errorf("Expected status in full render:\n%s", output)
	}
	if !strings.Contains(output, "Commands:") {
		t.Errorf("Expected help in full render:\n%s", output)
	}
}

func TestStatusStringEmptyChains(t *testing.T) {
	lvl := &level.Level{ID: 1, Values: [][]int{{1, 2}, {3, 4}}}
	g := game.NewGame(lvl, 0, 2, 2)
	activeChains := []int{}

	output := StatusString(g, 0, activeChains)
	if output != "" {
		t.Errorf("Expected empty status for no chains, got %q", output)
	}
}
