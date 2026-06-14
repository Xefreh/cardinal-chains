package levelloader

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"cardinal-chains/level"

	"gopkg.in/yaml.v3"
)

// yamlLevels mirrors the structure of levels.yml.
type yamlLevels struct {
	Levels map[int][]string `yaml:"levels"`
}

// parseRow parses a space-separated string of integers (e.g. "-1 1 1 1")
// into a slice of ints.
func parseRow(s string) ([]int, error) {
	fields := strings.Fields(s)
	values := make([]int, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q in row %q: %w", f, s, err)
		}
		values = append(values, v)
	}
	return values, nil
}

// ReadYAMLFile reads and parses the YAML level file, returning an ordered
// Levels collection sorted by level number.
func ReadYAMLFile(filename string) (*level.Levels, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	var yl yamlLevels
	if err := yaml.Unmarshal(data, &yl); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	levels := level.NewLevels()

	for id, rows := range yl.Levels {
		values := make([][]int, 0, len(rows))
		for _, row := range rows {
			parsed, err := parseRow(row)
			if err != nil {
				return nil, err
			}
			values = append(values, parsed)
		}
		levels.AddLevel(id, values)
	}

	levels.SortByID()

	return levels, nil
}
