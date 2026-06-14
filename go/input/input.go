package input

import "unicode"

// Command represents a player action parsed from a key press.
type Command int

const (
	Invalid Command = iota
	MoveNorth
	MoveSouth
	MoveEast
	MoveWest
	Cancel
	Erase
	Restart
	CycleChain
	Quit
)

// IsMove returns true if the command is one of the four movement directions.
func (c Command) IsMove() bool {
	return c == MoveNorth || c == MoveSouth || c == MoveEast || c == MoveWest
}

// String returns a human-readable name for the command.
func (c Command) String() string {
	switch c {
	case MoveNorth:
		return "Move North"
	case MoveSouth:
		return "Move South"
	case MoveEast:
		return "Move East"
	case MoveWest:
		return "Move West"
	case Cancel:
		return "Cancel Last Move"
	case Erase:
		return "Erase Chain"
	case Restart:
		return "Restart Level"
	case CycleChain:
		return "Cycle Chain"
	case Quit:
		return "Quit"
	default:
		return "Invalid"
	}
}

// ParseChar converts a single character into a Command, mirroring the
// switch on toupper(input) in the C game_loop.c. This function is
// case-insensitive.
func ParseChar(ch rune) Command {
	switch unicode.ToUpper(ch) {
	case 'N':
		return MoveNorth
	case 'S':
		return MoveSouth
	case 'E':
		return MoveEast
	case 'W':
		return MoveWest
	case 'B':
		return Cancel
	case 'R':
		return Erase
	case 'X':
		return Restart
	case 'C':
		return CycleChain
	case 'Q':
		return Quit
	default:
		return Invalid
	}
}
