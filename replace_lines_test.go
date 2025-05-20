package patt_test

import (
	"bytes"
	"patt"
	"strings"
	"testing"
)


func TestReplaceMultiline(t *testing.T) {
	replacer, _ := patt.NewReplacer("<match>","<match>")
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "handles long lines",
			input:    strings.Repeat("a", 10000) + "\n",
			expected: strings.Repeat("a", 10000) + "\n",
		},
		{
			name:     "handles new lines correctly",
			input:    "line1\nline2\nline3\n",
			expected: "line1\nline2\nline3\n",
		},
		{
			name:     "handles missing new line EOF",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3\n",
		},
		{
			name:     "handles empty lines",
			input:    "line1\n\nline2\nline3\n\n",
			expected: "line1\nline2\nline3\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := bytes.NewReader([]byte(tt.input))
			output := &bytes.Buffer{}

			matched, err := patt.ReplaceLines(replacer, input, output)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !matched {
				t.Errorf("expected match but got no match")
			}
			if output.String() != tt.expected {
				t.Errorf("expected output %q but got %q", tt.expected, output.String())
			}
		})
	}
}

func TestReplaceMatchingLines(t *testing.T) {
	pattern := "one <name> three"
	template := "1 <name> 3"
	tests := []struct {
		name     string
		input    string
		expected string
		match    bool
	}{
		{
			name:  "no lines, no match",
			input: "",
		},
		{
			name:  "one line, no match",
			input: "four five six",
		},
		{
			name:     "one line, 1 match",
			input:    "one two three",
			expected: "1 two 3\n",
			match:    true,
		},
		{
			name:     "2 lines, 1 match",
			input:    "one two three\nfour five six",
			expected: "1 two 3\n",
			match:    true,
		},
		{
			name:     "2 lines, no match",
			input:    "one five six\nseven eight nine",
			expected: "",
		},
		{
			name:     "2 lines, 2 matches",
			input:    "one two three\none 2 three",
			expected: "1 two 3\n1 2 3\n",
			match:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := makeReplacer(t, pattern, template)
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer

			matched, err := patt.ReplaceLines(matcher, reader, &writer)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.match != matched {
				t.Errorf("Expected match to be %v but got %v", tt.match, matched)
			}
			if writer.String() != tt.expected {
				t.Errorf("expected output %q but got %q", tt.expected, writer.String())
			}
		})
	}
}