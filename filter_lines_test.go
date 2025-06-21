package patt_test

import (
	"bytes"
	"io"
	"os"
	"patt"
	"strings"
	"testing"
)

func TestPrintMatchingLines_Merged(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected string
		match    bool
	}{
		// From TestPrintMultiline (pattern: <_>)
		{
			name:     "handles long lines",
			pattern:  "<_>",
			input:    strings.Repeat("a", 10000) + "\n",
			expected: strings.Repeat("a", 10000) + "\n",
			match:    true,
		},
		{
			name:     "handles new lines correctly",
			pattern:  "<_>",
			input:    "line1\nline2\nline3\n",
			expected: "line1\nline2\nline3\n",
			match:    true,
		},
		{
			name:     "handles missing new line EOF",
			pattern:  "<_>",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3\n",
			match:    true,
		},
		{
			name:     "handles empty lines",
			pattern:  "<_>",
			input:    "line1\n\nline2\nline3\n\n",
			expected: "line1\nline2\nline3\n",
			match:    true,
		},
		// From TestPrintMatchingLines (pattern: one <_> three)
		{
			name:     "no lines, no match",
			pattern:  "one <_> three",
			input:    "",
			expected: "",
			match:    false,
		},
		{
			name:     "one line, no match",
			pattern:  "one <_> three",
			input:    "four five six",
			expected: "",
			match:    false,
		},
		{
			name:     "one line, 1 match",
			pattern:  "one <_> three",
			input:    "one two three",
			expected: "one two three\n",
			match:    true,
		},
		{
			name:     "2 lines, 1 match",
			pattern:  "one <_> three",
			input:    "one two three\nfour five six",
			expected: "one two three\n",
			match:    true,
		},
		{
			name:     "2 lines, no match",
			pattern:  "one <_> three",
			input:    "four five six\nseven eight nine",
			expected: "",
			match:    false,
		},
		{
			name:     "2 lines, 2 matches",
			pattern:  "one <_> three",
			input:    "one two three\none 2 three",
			expected: "one two three\none 2 three\n",
			match:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher := makeMatcher(t, tt.pattern)
			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer
			processor := patt.NewLineProcessor(reader, &writer, false)

			match, err := processor.ProcessLines(matcher)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.match != match {
				t.Errorf("Expected match to be %v but got %v", tt.match, match)
			}
			if writer.String() != tt.expected {
				t.Errorf("expected output %q but got %q", tt.expected, writer.String())
			}
		})
	}
}

func makeMatcher(t testing.TB, stringPattern string) patt.LineReplacer {
	t.Helper()
	matcher, err := patt.NewFilter(stringPattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return matcher
}

func TestParseApacheLogFile(t *testing.T) {
	filePath := "testdata/Apache_2k.log"

	input, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	var writer bytes.Buffer
	processor := patt.NewLineProcessor(input, &writer, false)
	matcher := makeMatcher(t, "[<_>] [error] <_>")

	match, err := processor.ProcessLines(matcher)
	if err != nil {
		t.Fatalf("error during matching: %v", err)
	}
	if !match {
		t.Errorf("no match")
	}
}

func BenchmarkParseLargeString(b *testing.B) {
	lines := `something once
something twice
something thrice
not matching once
someone once
`
	fileSize := 500 * 1024 * 1024 // 500 MB
	lineSize := len(lines)
	times := fileSize / lineSize

	var buffer bytes.Buffer
	for range times {
		buffer.WriteString(lines)
	}
	fileContent := buffer.Bytes()
	reader := bytes.NewReader(fileContent) // Reusable reader
	writer := io.Discard
	processor := patt.NewLineProcessor(reader, writer, false)

	matcher := makeMatcher(b, "something <_>")

	b.SetBytes(int64(len(fileContent)))

	b.ResetTimer()
	for b.Loop() {
		_, _ = reader.Seek(0, io.SeekStart)
		match, err := processor.ProcessLines(matcher)
		if err != nil {
			b.Fatalf("error during matching: %v", err)
		}
		if !match {
			b.Fatalf("no match")
		}
	}
}

func BenchmarkParseLargeFile(b *testing.B) {
	filePath := "testdata/Apache_500MB.log"
	reader, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		b.Fatalf("failed to open file: %v", err)
	}
	defer reader.Close()
	fileInfo, err := reader.Stat()
	if err != nil {
		b.Fatalf("failed to stat file: %v", err)
	}
	b.SetBytes(fileInfo.Size())
	matcher := makeMatcher(b, "[<_>] [error] <_>")
	writer := io.Discard
	processor := patt.NewLineProcessor(reader, writer, false)

	b.ResetTimer()
	for b.Loop() {
		_, _ = reader.Seek(0, io.SeekStart)
		match, err := processor.ProcessLines(matcher)
		if err != nil {
			b.Fatalf("error during matching: %v", err)
		}
		if !match {
			b.Fatalf("no match")
		}
	}
}

func BenchmarkParseMemoryLoadedFile(b *testing.B) {
	filePath := "testdata/Apache_500MB.log"
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		b.Fatalf("failed to read file: %v", err)
	}
	b.SetBytes(int64(len(fileContent)))

	matcher := makeMatcher(b, "[<_>] [error] <_>")

	writer := io.Discard
	reader := bytes.NewReader(fileContent)
	processor := patt.NewLineProcessor(reader, writer, false)
	b.ResetTimer()
	for b.Loop() {
		_, _ = reader.Seek(0, io.SeekStart)
		match, err := processor.ProcessLines(matcher)
		if err != nil {
			b.Fatalf("error during matching: %v", err)
		}
		if !match {
			b.Fatalf("no match")
		}
	}
}
