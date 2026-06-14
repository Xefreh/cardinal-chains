package main

import (
	"fmt"
	"os"

	"cardinal-chains/game"
	"cardinal-chains/gameloop"
	"cardinal-chains/levelloader"

	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <YAML file>\n", os.Args[0])
		os.Exit(1)
	}

	levels, err := levelloader.ReadYAMLFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading levels: %v\n", err)
		os.Exit(1)
	}

	if levels.Count() == 0 {
		fmt.Fprintf(os.Stderr, "No levels found in %s\n", os.Args[1])
		os.Exit(1)
	}

	gridRows := levels.GridRows()
	gridCols := levels.GridCols()

	g := game.NewGame(&levels.Items[0], 0, gridRows, gridCols)

	app := tview.NewApplication()
	gl := gameloop.NewGameLoop(g, levels, app)

	if err := gl.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
