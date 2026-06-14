package gameloop

import (
	"fmt"
	"strings"

	"cardinal-chains/game"
	"cardinal-chains/input"
	"cardinal-chains/level"
	"cardinal-chains/logger"
	"cardinal-chains/render"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// GameLoop is the controller that wires input → game rules → rendering
// each turn, mirroring the C game_loop.c play_game() function.
type GameLoop struct {
	game         *game.Game
	levels       *level.Levels
	app          *tview.Application
	gridView     *tview.TextView
	statusView   *tview.TextView
	helpView     *tview.TextView
	currentChain int
	activeChains []int
}

// NewGameLoop creates a new game loop controller with the given game state,
// levels collection, and tview application.
func NewGameLoop(g *game.Game, levels *level.Levels, app *tview.Application) *GameLoop {
	gl := &GameLoop{
		game:   g,
		levels: levels,
		app:    app,
	}

	gl.activeChains = make([]int, g.ChainCount())
	for i := range gl.activeChains {
		gl.activeChains[i] = i
	}

	gl.setupUI()
	return gl
}

// setupUI builds the tview layout: title, grid, status, and help panels.
func (gl *GameLoop) setupUI() {
	gl.gridView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetWrap(false)
	gl.gridView.SetBorder(true).SetTitle(" Cardinal Chains ").SetTitleAlign(tview.AlignLeft)

	gl.statusView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	gl.statusView.SetBorder(true).SetTitle(" Status ").SetTitleAlign(tview.AlignLeft)

	gl.helpView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)
	gl.helpView.SetBorder(true).SetTitle(" Commands ").SetTitleAlign(tview.AlignLeft)
	gl.helpView.SetText(render.HelpString())

	gl.refresh()
}

// Layout returns the root tview primitive for the game.
func (gl *GameLoop) Layout() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(gl.gridView, 0, 1, true).
		AddItem(gl.statusView, 4, 0, false).
		AddItem(gl.helpView, 8, 0, false)
	return flex
}

// refresh re-renders the grid and status panels from the current game state.
func (gl *GameLoop) refresh() {
	fmt.Fprintf(gl.gridView, "%s", render.GridString(gl.game, gl.currentChain, gl.activeChains))
	fmt.Fprintf(gl.statusView, "%s", render.StatusString(gl.game, gl.currentChain, gl.activeChains))
}

// showLevelComplete displays the level completion overlay.
func (gl *GameLoop) showLevelComplete() {
	gl.statusView.Clear()
	fmt.Fprintf(gl.statusView, "\n[green::b]Level completed![white:-:-]\n")
}

// handleCommand processes a parsed command against the game state.
// Returns true if the application should continue running, false if it
// should quit.
func (gl *GameLoop) handleCommand(cmd input.Command) bool {
	switch cmd {
	case input.Quit:
		gl.app.Stop()
		return false

	case input.MoveNorth, input.MoveSouth, input.MoveEast, input.MoveWest:
		dir := dirFromCommand(cmd)
		moved := gl.game.MoveChain(gl.currentChain, dir)
		if moved && gl.game.IsGameCompleted() {
			return gl.advanceLevel()
		}

	case input.Cancel:
		gl.game.CancelLastMove(gl.currentChain)

	case input.Erase:
		gl.game.EraseChain(gl.currentChain)

	case input.Restart:
		gl.game.RestartLevel()

	case input.CycleChain:
		if gl.game.ChainCount() > 1 {
			gl.currentChain = (gl.currentChain + 1) % gl.game.ChainCount()
		}

	case input.Invalid:
		logger.Debug.Println("invalid command")
	}

	gl.refresh()
	return true
}

// advanceLevel handles the transition to the next level after completion.
// Returns false if the game should stop (all levels done), true otherwise.
func (gl *GameLoop) advanceLevel() bool {
	gl.showLevelComplete()
	gl.game.CurrentLevel++

	if gl.game.CurrentLevel < gl.levels.Count() {
		gl.game.LoadNextLevel(gl.levels)
		gl.currentChain = 0
		gl.activeChains = make([]int, gl.game.ChainCount())
		for i := range gl.activeChains {
			gl.activeChains[i] = i
		}
		gl.refresh()
		return true
	}

	gl.allLevelsComplete()
	return false
}

// allLevelsComplete shows the congratulations message and stops the app.
func (gl *GameLoop) allLevelsComplete() {
	gl.gridView.Clear()
	msg := strings.Join([]string{
		"",
		"[green::b]Congratulations!",
		"You have completed all levels![white:-:-]",
	}, "\n")
	fmt.Fprintf(gl.gridView, "%s\n", msg)
}

// dirFromCommand maps a movement command to a game.Direction.
func dirFromCommand(cmd input.Command) game.Direction {
	switch cmd {
	case input.MoveNorth:
		return game.North
	case input.MoveSouth:
		return game.South
	case input.MoveEast:
		return game.East
	case input.MoveWest:
		return game.West
	default:
		return game.North
	}
}

// HandleKey is the input capture callback for tview. It converts the key
// event into a command and processes it.
func (gl *GameLoop) HandleKey(event *tcell.EventKey) *tcell.EventKey {
	ch := event.Rune()
	if ch == 0 {
		return event
	}

	cmd := input.ParseChar(ch)
	if cmd == input.Invalid {
		return event
	}

	gl.handleCommand(cmd)
	return nil
}

// Run starts the game loop by setting the root primitive and input capture,
// then running the tview application.
func (gl *GameLoop) Run() error {
	root := gl.Layout()
	gl.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return gl.HandleKey(event)
	})
	gl.app.SetRoot(root, true)
	return gl.app.Run()
}
