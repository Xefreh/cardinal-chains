package level

import "sort"

// Level represents a single puzzle grid.
//
// ID is the level number as read from the YAML file. Values is a jagged
// 2-D grid of integers. Each cell value has special meaning:
//
//   - -1 — an anchor cell. Anchors always come in pairs and form the
//     start/end of a chain. Each -1 becomes the seed of one chain.
//   - 0  — an empty/blocked cell that cannot be entered by any chain.
//   - A positive integer — a fillable cell. Every positive cell must be
//     covered by some chain for the level to be considered complete.
type Level struct {
	ID     int
	Values [][]int
}

// Levels is an ordered collection of Level objects.
type Levels struct {
	Items []Level
}

// NewLevels creates an empty Levels collection.
func NewLevels() *Levels {
	return &Levels{Items: []Level{}}
}

// AddLevel appends a new level to the collection.
func (ls *Levels) AddLevel(id int, values [][]int) {
	ls.Items = append(ls.Items, Level{ID: id, Values: values})
}

// Count returns the number of levels in the collection.
func (ls *Levels) Count() int {
	return len(ls.Items)
}

// Rows returns the number of rows in this level's grid.
func (l *Level) Rows() int {
	return len(l.Values)
}

// Cols returns the maximum column count across all rows, mirroring the
// C implementation which uses the widest row as the grid width.
func (l *Level) Cols() int {
	maxCols := 0
	for _, row := range l.Values {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	return maxCols
}

// CountChains counts the number of anchor cells (value == -1). Each anchor
// cell is the seed of one chain.
func (l *Level) CountChains() int {
	count := 0
	for _, row := range l.Values {
		for _, v := range row {
			if v == -1 {
				count++
			}
		}
	}
	return count
}

// GridRows returns the maximum row count across all levels (the global
// grid height), mirroring the C implementation in main.c.
func (ls *Levels) GridRows() int {
	maxRows := 0
	for _, l := range ls.Items {
		if l.Rows() > maxRows {
			maxRows = l.Rows()
		}
	}
	return maxRows
}

// GridCols returns the maximum column count across all rows across all
// levels (the global grid width), mirroring the C implementation in main.c.
func (ls *Levels) GridCols() int {
	maxCols := 0
	for _, l := range ls.Items {
		if c := l.Cols(); c > maxCols {
			maxCols = c
		}
	}
	return maxCols
}

// SortByID sorts the levels in ascending order by their ID. This ensures
// deterministic iteration order regardless of how the underlying YAML map
// was parsed.
func (ls *Levels) SortByID() {
	sort.Slice(ls.Items, func(i, j int) bool {
		return ls.Items[i].ID < ls.Items[j].ID
	})
}

// CellValue safely retrieves the value at (row, col). Returns 0 for
// out-of-bounds positions (treated as blocked cells).
func (l *Level) CellValue(row, col int) int {
	if row < 0 || row >= len(l.Values) {
		return 0
	}
	if col < 0 || col >= len(l.Values[row]) {
		return 0
	}
	return l.Values[row][col]
}
