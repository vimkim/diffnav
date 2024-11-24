package ripgrep

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
)

// Structs to match ripgrep's JSON output
type MatchObject struct {
	Path       string     `json:"path"`
	Lines      string     `json:"lines"`
	Submatches []Submatch `json:"submatches"`
	LineNumber int        `json:"line_number"`
}

type Submatch struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

// Struct to parse raw ripgrep JSON
type RipgrepRaw struct {
	Type string `json:"type"`
	Data struct {
		Path struct {
			Text string `json:"text"`
		} `json:"path"`
		Lines struct {
			Text string `json:"text"`
		} `json:"lines"`
		Submatches []struct {
			Match struct {
				Text string `json:"text"`
			} `json:"match"`
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"submatches"`
		LineNumber int `json:"line_number"`
	} `json:"data"`
}

func (m MatchObject) GetFileName() string {
	return filepath.Base(m.Path)
}

func (m MatchObject) Dir() string {
	return filepath.Dir(m.Path)
}

// Updated Parse function to handle multiple JSON objects in input
func Parse(input io.Reader) ([]*MatchObject, error) {
	var results []*MatchObject

	scanner := bufio.NewScanner(input) // Read input line by line
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		var raw RipgrepRaw
		err := json.Unmarshal([]byte(line), &raw)
		if err != nil {
			return nil, fmt.Errorf("error parsing line: %w", err)
		}

		// Only process "match" type entries
		if raw.Type != "match" {
			continue
		}

		// Map raw data to MatchObject
		match := &MatchObject{
			Path:       raw.Data.Path.Text,
			LineNumber: raw.Data.LineNumber,
			Lines:      raw.Data.Lines.Text,
		}

		// Map submatches
		for _, sub := range raw.Data.Submatches {
			match.Submatches = append(match.Submatches, Submatch{
				Text:  sub.Match.Text,
				Start: sub.Start,
				End:   sub.End,
			})
		}

		results = append(results, match)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input: %w", err)
	}

	return results, nil
}
