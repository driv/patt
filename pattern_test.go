package patt_test

import (
	"bytes"
	"patt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMatcher(t *testing.T) {
	tests := []struct {
		name          string
		stringPattern string
		wantErr       bool
		wantNil       bool
	}{
		{
			name:          "Correct stringPattern",
			stringPattern: "something <_> something else",
			wantErr:       false,
			wantNil:       false,
		},
		{
			name:          "Incorrect stringPattern",
			stringPattern: "something <_><_> something else",
			wantErr:       true,
			wantNil:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := patt.NewMatcher(tt.stringPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatcher() error = %v, expected err: %v", err, tt.wantErr)
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewMatcher() = %v, expected nil: %v", got, tt.wantNil)
			}
		})
	}
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
