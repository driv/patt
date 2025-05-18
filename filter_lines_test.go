package patt_test

import (
	"bytes"
	"io"
	"os"
	"patt"
	"regexp"
	"strings"
	"testing"
)

func TestPrintMultiline(t *testing.T) {
	matcher, _ := patt.NewMatcher("<_>")
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

			matched, err := patt.PrintMatchingLines(matcher, input, output)

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

func TestPrintMatchingLines(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected string
		match  bool
	}{
		{
			name:    "no lines, no match",
			pattern: "<_>",
			input:   "",
		},
		{
			name:    "one line, no match",
			pattern: "something <_>",
			input:   "wrong stringPattern",
		},
		{
			name:     "one line, 1 match",
			pattern:  "something <_>",
			input:    "something stringPattern",
			expected: "something stringPattern\n",
			match: true,
		},
		{
			name:     "2 lines, 1 match",
			pattern:  "something <_>",
			input:    "wrong stringPattern\nsomething stringPattern",
			expected: "something stringPattern\n",
			match: true,
		},
		{
			name:     "2 lines, no match",
			pattern:  "something <_>",
			input:    "wrong stringPattern\nnon-matching stringPattern",
			expected: "",
		},
		{
			name:     "2 lines, 2 matches",
			pattern:  "something <_>Pattern",
			input:    "something oncePattern\nsomething twicePattern",
			expected: "something oncePattern\nsomething twicePattern\n",
			match: true,
		},
		{
			name:     "spaces are not special",
			pattern:  "something <_>Pattern",
			input:    "something  oncePattern\nsomething  twicePattern",
			expected: "something  oncePattern\nsomething  twicePattern\n",
			match: true,
		},
		{
			name:     "parentheses and stuff",
			pattern:  "[<_>] [error] <_>",
			input:    "[01:01:01] [error] some error message",
			expected: "[01:01:01] [error] some error message\n",
			match: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := makeMatcher(t, tt.pattern)
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer

			matched, err := patt.PrintMatchingLines(matcher, reader, &writer)

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

func makeMatcher(t testing.TB, stringPattern string) patt.LinesMatcher {
	t.Helper()
	matcher, err := patt.NewMatcher(stringPattern)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	return matcher
}

func BenchmarkMatchLines(b *testing.B) {
	readerContent := strings.Repeat(`something once
something twice
something thrice
not matching once
someone once
`, 10)
	filter := makeMatcher(b, "somet<_> <_>")
	b.Run("pattern matcher", func(b *testing.B) {

		for b.Loop() {
			reader := strings.NewReader(readerContent)
			var writer bytes.Buffer
			_, err := patt.PrintMatchingLines(filter, reader, &writer)
			if err != nil {
				return
			}
		}
	})

	regex := `somet.+ .+`
	re := regexp.MustCompile(regex)
	b.Run("regex matcher", func(b *testing.B) {
		for b.Loop() {
			reader := strings.NewReader(readerContent)
			var writer bytes.Buffer
			_, err := patt.PrintMatchingLines(re, reader, &writer)
			if err != nil {
				return
			}
		}
	})
}

func BenchmarkParseLargeFile(b *testing.B) {
	lines := strings.Repeat(`something once
something twice
something thrice
not matching once
someone once
`, 10)
	fileSize := 500 * 1024 * 1024 // 500 MB
	times := fileSize / len(lines)
	var builder strings.Builder
	for i := 0; i < times; i++ {
		builder.WriteString(lines)
	}
	fileContent := builder.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(fileContent)
		var writer bytes.Buffer
		matcher := makeMatcher(b, "something <_>")
		_, err := patt.PrintMatchingLines(matcher, reader, &writer)
		if err != nil {
			b.Fatalf("error during matching: %v", err)
		}
	}
}

func TestParseApacheLogFile(t *testing.T) {
	filePath := "test_files/Apache_2k.log"

	input, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var writer bytes.Buffer
	matcher := makeMatcher(t, "[<_>] [error] <_>")
	match, err := patt.PrintMatchingLines(matcher, input, &writer)
	if err != nil {
		t.Fatalf("error during matching: %v", err)
	}
	if !match {
		t.Errorf("no match")
	}
}

func BenchmarkParseApacheLogFileHuge(b *testing.B) {
	filePath := "test_files/Apache_500MB.log"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input, err := os.OpenFile(filePath, os.O_RDONLY, 0)
		if err != nil {
			b.Fatalf("failed to read file: %v", err)
		}
		matcher := makeMatcher(b, "[<_>] [error] <_>")
		match, err := patt.PrintMatchingLines(matcher, input, io.Discard)
		if err != nil {
			b.Fatalf("error during matching: %v", err)
		}
		if match != true {
			b.Fatalf("no match")
		}
	}
}
