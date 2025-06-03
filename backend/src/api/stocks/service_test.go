package stocks

import (
	"testing"
)

func TestNormalizeRating(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Buy", 9},
		{"buy", 9},
		{"SELL", 1},
		{"Neutral", 5},
		{"Strong-Buy", 10},
		{"unknown", 0},
		{"", 0},
	}

	for _, tt := range tests {
		got := normalizeRating(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeRating(%q) = %v; want %v", tt.input, got, tt.expected)
		}
	}
}
