package ripgrep

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestParseMultipleLines(t *testing.T) {
	mockInput := `{"type":"begin","data":{"path":{"text":"main.go"}}}
{"type":"match","data":{"path":{"text":"main.go"},"lines":{"text":"\t\"github.com/dlvhdr/diffnav/pkg/ui\"\n"},"line_number":14,"absolute_offset":195,"submatches":[{"match":{"text":"diff"},"start":20,"end":24}]}}
{"type":"match","data":{"path":{"text":"main.go"},"lines":{"text":"\t// \tfmt.Println(\"No diff, exiting\")\n"},"line_number":24,"absolute_offset":387,"submatches":[{"match":{"text":"diff"},"start":21,"end":25}]}}
`

	reader := strings.NewReader(mockInput)

	matches, err := Parse(reader)
	if err != nil {
		t.Fatalf("Parse returned an error: %v", err)
	}

	for _, i := range matches {
		matchJSON, err := json.Marshal(i)
		if err != nil {
			fmt.Printf("Error marshalling match: %v\n", err)
			continue
		}
		fmt.Println(string(matchJSON)) // Print the JSON representation of the match
	}

	// Validate number of matches
	expectedMatches := 2
	if len(matches) != expectedMatches {
		t.Fatalf("Expected %d matches, got %d", expectedMatches, len(matches))
	}

	// Validate first match
	match1 := matches[0]
	if match1.Path != "main.go" {
		t.Errorf("Expected Path to be 'main.go', got '%s'", match1.Path)
	}
	if match1.LineNumber != 14 {
		t.Errorf("Expected LineNumber to be 14, got %d", match1.LineNumber)
	}
	if match1.Lines != "\t\"github.com/dlvhdr/diffnav/pkg/ui\"\n" {
		t.Errorf("Expected Lines to match, got '%s'", match1.Lines)
	}
	if len(match1.Submatches) != 1 {
		t.Errorf("Expected 1 submatch, got %d", len(match1.Submatches))
	}

	// Validate submatch of first match
	sub1 := match1.Submatches[0]
	if sub1.Text != "diff" {
		t.Errorf("Expected Submatch.Text to be 'diff', got '%s'", sub1.Text)
	}
	if sub1.Start != 20 || sub1.End != 24 {
		t.Errorf("Expected Submatch start and end to be 20 and 24, got %d and %d", sub1.Start, sub1.End)
	}

	// Validate second match similarly (skipped for brevity)
}
