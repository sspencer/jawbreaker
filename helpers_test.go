package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestExtractNumberFromPieceId(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"valid zero", "piece0", 0},
		{"valid positive", "piece123", 123},
		{"invalid prefix", "foo123", -1},
		{"prefix only", "piece", -1},
		{"non-number suffix", "pieceabc", -1},
		{"negative number", "piece-5", -1}, // Assuming only non-negative allowed
		{"empty string", "", -1},
		{"number only", "123", -1},
		{"prefix with space", "piece 1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractNumberFromPieceId(tt.input); got != tt.want {
				t.Errorf("extractNumberFromPieceId(%q) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestGameToHTML(t *testing.T) {
	// Create a simple 2x2 game
	g := &Game{
		pieces: []Piece{Red, Blue, Green, White},
		rows:   2,
		cols:   2,
		score:  0,
	}

	html := gameToHTML(g)

	// Basic structural checks
	if !strings.HasPrefix(html, "<div id=\"game\"") {
		t.Errorf("gameToHTML() output does not start with expected game div")
	}
	if !strings.Contains(html, "data-on-click=\"@post('move/'+evt.target.id)\"") {
		t.Errorf("gameToHTML() output missing expected data-on-click attribute")
	}
	if !strings.HasSuffix(html, "</div>") {
		t.Errorf("gameToHTML() output does not end with expected closing div")
	}

	// Check for expected pieces
	expectedPiecesHTML := []string{
		"<div id=\"piece0\" class=\"red\"></div>",
		"<div id=\"piece1\" class=\"blue\"></div>",
		"<div id=\"piece2\" class=\"green\"></div>",
		"<div id=\"piece3\" class=\"white\"></div>",
	}

	for _, pieceHTML := range expectedPiecesHTML {
		if !strings.Contains(html, pieceHTML) {
			t.Errorf("gameToHTML() output missing expected piece HTML: %s", pieceHTML)
		}
	}

	// Check count of divs roughly matches
	// Add 1 for the outer game div
	expectedDivCount := len(g.pieces) + 1
	actualDivCount := strings.Count(html, "<div")
	if actualDivCount != expectedDivCount {
		t.Errorf("gameToHTML() expected %d opening divs, found %d", expectedDivCount, actualDivCount)
	}
}

func TestScoreData_Serialization(t *testing.T) {
	tests := []struct {
		name       string
		input      ScoreData
		wantString string
	}{
		{"zero scores", ScoreData{LastScore: 0, BestScore: 0}, "0|0"},
		{"positive scores", ScoreData{LastScore: 123, BestScore: 456}, "123|456"},
		{"mixed scores", ScoreData{LastScore: 0, BestScore: 999}, "0|999"},
		// Add tests for potentially large scores if applicable
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Serialize_%s", tt.name), func(t *testing.T) {
			if got := tt.input.serialize(); got != tt.wantString {
				t.Errorf("ScoreData.serialize() for %+v = %q; want %q", tt.input, got, tt.wantString)
			}
		})
	}

	// --- Test Deserialization ---
	deserializeTests := []struct {
		name        string
		inputString string
		wantData    ScoreData
		initialData ScoreData // What the struct holds before deserialization
	}{
		{"valid scores", "123|456", ScoreData{LastScore: 123, BestScore: 456}, ScoreData{}},
		{"zero scores", "0|0", ScoreData{LastScore: 0, BestScore: 0}, ScoreData{}},
		{"invalid format (no pipe)", "123456", ScoreData{LastScore: 50, BestScore: 100}, ScoreData{LastScore: 50, BestScore: 100}},             // Should not change
		{"invalid format (too many pipes)", "1|2|3", ScoreData{LastScore: 50, BestScore: 100}, ScoreData{LastScore: 50, BestScore: 100}},       // Should not change
		{"invalid format (non-numeric first)", "abc|456", ScoreData{LastScore: 50, BestScore: 100}, ScoreData{LastScore: 50, BestScore: 100}},  // Should not change
		{"invalid format (non-numeric second)", "123|def", ScoreData{LastScore: 50, BestScore: 100}, ScoreData{LastScore: 50, BestScore: 100}}, // Should not change
		{"empty string", "", ScoreData{LastScore: 50, BestScore: 100}, ScoreData{LastScore: 50, BestScore: 100}},                               // Should not change
	}

	for _, tt := range deserializeTests {
		t.Run(fmt.Sprintf("Deserialize_%s", tt.name), func(t *testing.T) {
			// Start with initial data
			gotData := tt.initialData
			deserializeScoreData(&gotData, tt.inputString)

			if gotData.LastScore != tt.wantData.LastScore || gotData.BestScore != tt.wantData.BestScore {
				t.Errorf("deserializeScoreData(%q) on initial %+v resulted in %+v; want %+v", tt.inputString, tt.initialData, gotData, tt.wantData)
			}
		})
	}
}
