package level

import "testing"

func TestNewLevels(t *testing.T) {
	ls := NewLevels()
	if ls.Count() != 0 {
		t.Errorf("NewLevels().Count() = %d, want 0", ls.Count())
	}
	if len(ls.Items) != 0 {
		t.Errorf("NewLevels().Items len = %d, want 0", len(ls.Items))
	}
}

func TestAddLevel(t *testing.T) {
	ls := NewLevels()
	ls.AddLevel(1, [][]int{{1, 2}, {3, 4}})
	ls.AddLevel(2, [][]int{{5, 6}})

	if ls.Count() != 2 {
		t.Fatalf("Count() = %d, want 2", ls.Count())
	}
	if ls.Items[0].ID != 1 {
		t.Errorf("Items[0].ID = %d, want 1", ls.Items[0].ID)
	}
	if ls.Items[1].ID != 2 {
		t.Errorf("Items[1].ID = %d, want 2", ls.Items[1].ID)
	}
}

func TestLevelRows(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}}
	if got := l.Rows(); got != 3 {
		t.Errorf("Rows() = %d, want 3", got)
	}
}

func TestLevelRowsEmpty(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{}}
	if got := l.Rows(); got != 0 {
		t.Errorf("Rows() = %d, want 0", got)
	}
}

func TestLevelCols(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{{1, 2, 3}, {4, 5}, {7}}}
	if got := l.Cols(); got != 3 {
		t.Errorf("Cols() = %d, want 3", got)
	}
}

func TestLevelColsUniform(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}}}
	if got := l.Cols(); got != 4 {
		t.Errorf("Cols() = %d, want 4", got)
	}
}

func TestCountChains(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{
		{-1, 1, 1, 1},
		{0, 0, 0, 0},
		{1, 1, 1, -1},
	}}
	if got := l.CountChains(); got != 2 {
		t.Errorf("CountChains() = %d, want 2", got)
	}
}

func TestCountChainsNoAnchors(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}}
	if got := l.CountChains(); got != 0 {
		t.Errorf("CountChains() = %d, want 0", got)
	}
}

func TestCountChainsAllAnchors(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{{-1, -1}, {-1, -1}}}
	if got := l.CountChains(); got != 4 {
		t.Errorf("CountChains() = %d, want 4", got)
	}
}

func TestCellValue(t *testing.T) {
	l := Level{ID: 1, Values: [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}}
	tests := []struct {
		name string
		row  int
		col  int
		want int
	}{
		{"in bounds (0,0)", 0, 0, 1},
		{"in bounds (1,2)", 1, 2, 6},
		{"in bounds (0,2)", 0, 2, 3},
		{"out of bounds row negative", -1, 0, 0},
		{"out of bounds row too large", 2, 0, 0},
		{"out of bounds col negative", 0, -1, 0},
		{"out of bounds col too large", 0, 3, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := l.CellValue(tt.row, tt.col); got != tt.want {
				t.Errorf("CellValue(%d, %d) = %d, want %d", tt.row, tt.col, got, tt.want)
			}
		})
	}
}

func TestLevelsGridRows(t *testing.T) {
	ls := NewLevels()
	ls.AddLevel(1, [][]int{{1, 2}, {3, 4}, {5, 6}})
	ls.AddLevel(2, [][]int{{1, 2}})
	ls.AddLevel(3, [][]int{{1}, {2}, {3}, {4}, {5}})

	if got := ls.GridRows(); got != 5 {
		t.Errorf("GridRows() = %d, want 5", got)
	}
}

func TestLevelsGridCols(t *testing.T) {
	ls := NewLevels()
	ls.AddLevel(1, [][]int{{1, 2, 3}, {3, 4}})
	ls.AddLevel(2, [][]int{{1, 2}})
	ls.AddLevel(3, [][]int{{1, 2, 3, 4, 5, 6}})

	if got := ls.GridCols(); got != 6 {
		t.Errorf("GridCols() = %d, want 6", got)
	}
}

func TestSortByID(t *testing.T) {
	ls := NewLevels()
	ls.AddLevel(3, [][]int{{1}})
	ls.AddLevel(1, [][]int{{1}})
	ls.AddLevel(2, [][]int{{1}})

	ls.SortByID()

	if ls.Items[0].ID != 1 || ls.Items[1].ID != 2 || ls.Items[2].ID != 3 {
		t.Errorf("SortByID order = %d, %d, %d; want 1, 2, 3",
			ls.Items[0].ID, ls.Items[1].ID, ls.Items[2].ID)
	}
}
