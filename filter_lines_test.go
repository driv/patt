package patt_test

import (
	"bytes"
	"io"
	"os"
	"patt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMatcher struct {
	matchFunc func([]byte) bool
}

func (m mockMatcher) Match(b []byte) bool {
	return m.matchFunc(b)
}

func TestMatchLines_LongLines(t *testing.T) {
	t.Run("handles long lines", func(t *testing.T) {
		longLine := strings.Repeat("a", 10000) + "\n"
		input := bytes.NewReader([]byte(longLine))
		output := &bytes.Buffer{}
		matcher := mockMatcher{
			matchFunc: func(b []byte) bool {
				return true
			},
		}

		matched, err := patt.PrintMatchingLines(matcher, input, output)

		assert.NoError(t, err)
		assert.True(t, matched)
		assert.Equal(t, longLine, output.String())
	})
}

func TestMatchLines_NewLines(t *testing.T) {
	t.Run("handles new lines correctly", func(t *testing.T) {
		input := bytes.NewReader([]byte("line1\nline2\nline3\n"))
		output := &bytes.Buffer{}
		matcher := mockMatcher{
			matchFunc: func(b []byte) bool {
				return true
			},
		}

		matched, err := patt.PrintMatchingLines(matcher, input, output)

		assert.NoError(t, err)
		assert.True(t, matched)
		assert.Equal(t, "line1\nline2\nline3\n", output.String())
	})
	t.Run("handles missing new line EOF", func(t *testing.T) {
		input := bytes.NewReader([]byte("line1\nline2\nline3"))
		output := &bytes.Buffer{}
		matcher := mockMatcher{
			matchFunc: func(b []byte) bool {
				return true
			},
		}

		matched, err := patt.PrintMatchingLines(matcher, input, output)

		assert.NoError(t, err)
		assert.True(t, matched)
		assert.Equal(t, "line1\nline2\nline3\n", output.String())
	})
	t.Run("handles empty lines", func(t *testing.T) {
		input := bytes.NewReader([]byte("line1\n\nline2\nline3\n\n"))
		output := &bytes.Buffer{}
		matcher := mockMatcher{
			matchFunc: func(b []byte) bool {
				return true
			},
		}

		matched, err := patt.PrintMatchingLines(matcher, input, output)

		assert.NoError(t, err)
		assert.True(t, matched)
		assert.Equal(t, "line1\n\nline2\nline3\n\n", output.String())
	})
}

func TestMatchesMultiple(t *testing.T) {
	t.Run("no lines, no match", func(t *testing.T) {
		matcher := makeMatcher(t, "<_>")
		reader := strings.NewReader("")
		var writer bytes.Buffer

		matched, err := patt.PrintMatchingLines(matcher, reader, &writer)

		assertNoMatch(t, err, matched, writer)
	})
	t.Run("one line, no match", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>")
		reader := strings.NewReader("wrong stringPattern")
		var writer bytes.Buffer

		matched, err := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.NoError(t, err, "MatchLines() should not return error")
		assert.False(t, matched, "MatchLines() should return false")
		assert.Empty(t, writer.Bytes(), "MatchLines() should not write")
	})
	t.Run("one line, 1 match", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>")
		reader := strings.NewReader("something stringPattern")
		var writer bytes.Buffer

		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.True(t, matched, "MatchLines() should return true")
		assert.NotEmpty(t, writer.Bytes(), "MatchLines() should have written to writer")
		assert.Equal(t, `something stringPattern
`, writer.String())
	})
	t.Run("2 lines, 1 match", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>")
		reader := strings.NewReader("wrong stringPattern\nsomething stringPattern")
		var writer bytes.Buffer

		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.True(t, matched, "MatchLines() should return true")
		assert.NotEmpty(t, writer.Bytes(), "MatchLines() should have written to writer")
		assert.Equal(t, `something stringPattern
`, writer.String())
	})
	t.Run("2 lines, no match", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>")
		reader := strings.NewReader(`wrong stringPattern
non-matching stringPattern`)
		var writer bytes.Buffer

		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.False(t, matched, "MatchLines()")
		assert.Empty(t, writer.Bytes(), "MatchLines() should not have written to writer")
	})
	t.Run("2 lines, 2 matches", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>Pattern")
		reader := strings.NewReader(`something oncePattern
something twicePattern`)
		var writer bytes.Buffer

		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.True(t, matched, "MatchLines() should return true")
		assert.NotEmpty(t, writer.Bytes(), "MatchLines() should have written to writer")
		assert.Equal(t, `something oncePattern
something twicePattern
`, writer.String())
	})
	t.Run("spaces are not special", func(t *testing.T) {
		matcher := makeMatcher(t, "something <_>Pattern")
		reader := strings.NewReader(`something  oncePattern
something  twicePattern`)
		var writer bytes.Buffer
		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.True(t, matched, "MatchLines() should return true")
		assert.NotEmpty(t, writer.Bytes(), "MatchLines() should have written to writer")
		assert.Equal(t, `something  oncePattern
something  twicePattern
`, writer.String())
	})
	t.Run("parentheses and stuff", func(t *testing.T) {
		matcher := makeMatcher(t, "[<_>] [error] <_>")
		reader := strings.NewReader(`[01:01:01] [error] some error message`)
		var writer bytes.Buffer
		matched, _ := patt.PrintMatchingLines(matcher, reader, &writer)

		assert.True(t, matched, "MatchLines() should return true")
		assert.NotEmpty(t, writer.Bytes(), "MatchLines() should have written to writer")
		assert.Equal(t, `[01:01:01] [error] some error message
`, writer.String())
	})
}

func assertNoMatch(t *testing.T, err error, matched bool, writer bytes.Buffer) {
	assert.NoError(t, err, "MatchLines() should not return error")
	assert.False(t, matched, "MatchLines() should return false")
	assert.Empty(t, writer.Bytes(), "MatchLines() should not write")
}

func makeMatcher(t testing.TB, stringPattern string) patt.LinesMatcher {
	t.Helper()
	matcher, err := patt.NewMatcher(stringPattern)
	assert.NoError(t, err, "NewMatcher() should not return an error")
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

		for i := 0; i < b.N; i++ {
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
		for i := 0; i < b.N; i++ {
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