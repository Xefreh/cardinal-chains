package levelloader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTestYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "levels.yml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test YAML: %v", err)
	}
	return path
}

func TestReadYAMLFile(t *testing.T) {
	path := writeTestYAML(t, `levels:
  1:
    - "-1 1 1 1"
    - "0 0 0 0"
    - "1 1 1 -1"
`)

	levels, err := ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile error: %v", err)
	}

	if levels.Count() != 1 {
		t.Fatalf("Count() = %d, want 1", levels.Count())
	}

	lvl := levels.Items[0]
	if lvl.ID != 1 {
		t.Errorf("ID = %d, want 1", lvl.ID)
	}
	if lvl.Rows() != 3 {
		t.Errorf("Rows() = %d, want 3", lvl.Rows())
	}
	if lvl.Cols() != 4 {
		t.Errorf("Cols() = %d, want 4", lvl.Cols())
	}
	if lvl.CountChains() != 2 {
		t.Errorf("CountChains() = %d, want 2", lvl.CountChains())
	}

	want := [][]int{{-1, 1, 1, 1}, {0, 0, 0, 0}, {1, 1, 1, -1}}
	for i, row := range want {
		for j, v := range row {
			if got := lvl.Values[i][j]; got != v {
				t.Errorf("Values[%d][%d] = %d, want %d", i, j, got, v)
			}
		}
	}
}

func TestReadYAMLFileMultipleLevels(t *testing.T) {
	path := writeTestYAML(t, `levels:
  2:
    - "-1 2"
    - "0 2"
  1:
    - "-1 1 1"
    - "1 1 -1"
`)

	levels, err := ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile error: %v", err)
	}

	if levels.Count() != 2 {
		t.Fatalf("Count() = %d, want 2", levels.Count())
	}

	if levels.Items[0].ID != 1 {
		t.Errorf("Items[0].ID = %d, want 1 (should be sorted)", levels.Items[0].ID)
	}
	if levels.Items[1].ID != 2 {
		t.Errorf("Items[1].ID = %d, want 2 (should be sorted)", levels.Items[1].ID)
	}
}

func TestReadYAMLFileNegativeNumbers(t *testing.T) {
	path := writeTestYAML(t, `levels:
  1:
    - "-1 -1"
    - "0 0"
`)

	levels, err := ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile error: %v", err)
	}

	lvl := levels.Items[0]
	if lvl.Values[0][0] != -1 || lvl.Values[0][1] != -1 {
		t.Errorf("Negative values not parsed correctly: %v", lvl.Values[0])
	}
}

func TestReadYAMLFileNonexistentFile(t *testing.T) {
	_, err := ReadYAMLFile("/nonexistent/path/to/file.yml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestReadYAMLFileInvalidYAML(t *testing.T) {
	path := writeTestYAML(t, `this is not: [valid yaml`)

	_, err := ReadYAMLFile(path)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestReadYAMLFileEmptyFile(t *testing.T) {
	path := writeTestYAML(t, ``)

	levels, err := ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile error: %v", err)
	}
	if levels.Count() != 0 {
		t.Errorf("Count() = %d, want 0 for empty file", levels.Count())
	}
}

func TestParseRow(t *testing.T) {
	tests := []struct {
		input   string
		want    []int
		wantErr bool
	}{
		{"-1 1 1 1", []int{-1, 1, 1, 1}, false},
		{"0 0 0 0", []int{0, 0, 0, 0}, false},
		{"1 2 3 4 5", []int{1, 2, 3, 4, 5}, false},
		{"-1", []int{-1}, false},
		{"", []int{}, false},
		{"1 abc 3", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseRow(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("len = %d, want %d", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("[%d] = %d, want %d", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestReadYAMLFileExistingProjectLevel(t *testing.T) {
	path := filepath.Join("..", "..", "levels.yml")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("Skipping: %s not found", path)
	}

	levels, err := ReadYAMLFile(path)
	if err != nil {
		t.Fatalf("ReadYAMLFile error: %v", err)
	}
	if levels.Count() == 0 {
		t.Error("Expected at least one level in the project levels.yml")
	}
}
