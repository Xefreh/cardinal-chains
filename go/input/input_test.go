package input

import "testing"

func TestParseChar(t *testing.T) {
	tests := []struct {
		input rune
		want  Command
	}{
		{'N', MoveNorth},
		{'S', MoveSouth},
		{'E', MoveEast},
		{'W', MoveWest},
		{'B', Cancel},
		{'R', Erase},
		{'X', Restart},
		{'C', CycleChain},
		{'Q', Quit},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := ParseChar(tt.input); got != tt.want {
				t.Errorf("ParseChar('%c') = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCharLowercase(t *testing.T) {
	tests := []struct {
		input rune
		want  Command
	}{
		{'n', MoveNorth},
		{'s', MoveSouth},
		{'e', MoveEast},
		{'w', MoveWest},
		{'b', Cancel},
		{'r', Erase},
		{'x', Restart},
		{'c', CycleChain},
		{'q', Quit},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			if got := ParseChar(tt.input); got != tt.want {
				t.Errorf("ParseChar('%c') = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCharInvalid(t *testing.T) {
	invalidChars := []rune{'a', 'z', '1', '!', ' ', '\n', '\t', 'D', 'F', 'G', 'H'}
	for _, ch := range invalidChars {
		t.Run(string(ch), func(t *testing.T) {
			if got := ParseChar(ch); got != Invalid {
				t.Errorf("ParseChar('%c') = %v, want Invalid", ch, got)
			}
		})
	}
}

func TestCommandIsMove(t *testing.T) {
	moves := []Command{MoveNorth, MoveSouth, MoveEast, MoveWest}
	for _, cmd := range moves {
		if !cmd.IsMove() {
			t.Errorf("%v.IsMove() = false, want true", cmd)
		}
	}

	nonMoves := []Command{Invalid, Cancel, Erase, Restart, CycleChain, Quit}
	for _, cmd := range nonMoves {
		if cmd.IsMove() {
			t.Errorf("%v.IsMove() = true, want false", cmd)
		}
	}
}

func TestCommandString(t *testing.T) {
	tests := []struct {
		cmd  Command
		want string
	}{
		{MoveNorth, "Move North"},
		{MoveSouth, "Move South"},
		{MoveEast, "Move East"},
		{MoveWest, "Move West"},
		{Cancel, "Cancel Last Move"},
		{Erase, "Erase Chain"},
		{Restart, "Restart Level"},
		{CycleChain, "Cycle Chain"},
		{Quit, "Quit"},
		{Invalid, "Invalid"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.cmd.String(); got != tt.want {
				t.Errorf("%v.String() = %q, want %q", tt.cmd, got, tt.want)
			}
		})
	}
}
