package render

import (
	"fmt"
	"strings"

	"cardinal-chains/game"
)

// chainColors maps a chain color index to a tview color tag. These mirror
// the 6 ANSI colors used in the C render.c implementation:
//
//	red, green, yellow, blue, magenta, cyan
var chainColors = []string{
	"[red]",
	"[green]",
	"[yellow]",
	"[blue]",
	"[magenta]",
	"[cyan]",
}

// colorReset is the tview tag to reset color formatting.
const colorReset = "[-:]"

// chainColor returns the tview color tag for the given color index, cycling
// through the available colors.
func chainColor(idx int) string {
	return chainColors[idx%len(chainColors)]
}

// colorizedName returns a human-readable color name for a color index.
func colorizedName(idx int) string {
	names := []string{"Red", "Green", "Yellow", "Blue", "Magenta", "Cyan"}
	return names[idx%len(names)]
}

// chainIndexAt returns the color index (from activeChains) of the chain
// occupying cell (row, col), or (-1, false) if no chain is there.
func chainIndexAt(g *game.Game, row, col int, activeChains []int) (int, bool) {
	for k, ch := range g.Chains {
		for _, pos := range ch {
			if pos.X == row && pos.Y == col {
				if k < len(activeChains) {
					return activeChains[k], true
				}
				return k, true
			}
		}
	}
	return -1, false
}

// GridString renders the entire game grid as a string suitable for a tview
// TextView. Chains are colored using tview color tags.
func GridString(g *game.Game, currentChain int, activeChains []int) string {
	var sb strings.Builder

	for i := range g.Level.Values {
		for j := range g.Level.Values[i] {
			value := g.Level.Values[i][j]

			colorIdx, isChain := chainIndexAt(g, i, j, activeChains)
			if isChain {
				sb.WriteString(chainColor(colorIdx))
				if value == -1 {
					sb.WriteString("x ")
				} else {
					sb.WriteString(fmt.Sprintf("%d ", value))
				}
				sb.WriteString(colorReset)
			} else if value != 0 {
				sb.WriteString(fmt.Sprintf("%d ", value))
			} else {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// StatusString renders the status line showing the current position and
// current chain color, mirroring the C render.c output.
func StatusString(g *game.Game, currentChain int, activeChains []int) string {
	var sb strings.Builder

	if len(g.Chains) == 0 || currentChain < 0 || currentChain >= len(g.Chains) {
		return ""
	}

	chain := g.Chains[currentChain]
	lastPos := chain[len(chain)-1]

	sb.WriteString(fmt.Sprintf("Current position: row %d, column %d\n",
		lastPos.X+1, lastPos.Y+1))

	colorIdx := currentChain
	if currentChain < len(activeChains) {
		colorIdx = activeChains[currentChain]
	}

	sb.WriteString(fmt.Sprintf("Current chain color: %sChain %d%s (%s)\n",
		chainColor(colorIdx), colorIdx+1, colorReset, colorizedName(colorIdx)))

	return sb.String()
}

// HelpString returns the help text listing all available commands.
func HelpString() string {
	return strings.Join([]string{
		"Commands:",
		"  N/S/E/W - Move the current chain",
		"  B       - Cancel (undo) the previous move",
		"  R       - Erase the current chain back to its anchor",
		"  X       - Restart the whole level",
		"  C       - Cycle selection to the next chain",
		"  Q       - Quit the game",
	}, "\n")
}

// FullRender combines the grid, status, and help into a single string.
func FullRender(g *game.Game, currentChain int, activeChains []int) string {
	var sb strings.Builder
	sb.WriteString(GridString(g, currentChain, activeChains))
	sb.WriteString("\n")
	sb.WriteString(StatusString(g, currentChain, activeChains))
	sb.WriteString("\n")
	sb.WriteString(HelpString())
	return sb.String()
}
