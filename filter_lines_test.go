package patt

import (
	"bytes"
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

		matched, err := PrintMatchingLines(matcher, input, output)

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

		matched, err := PrintMatchingLines(matcher, input, output)

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

		matched, err := PrintMatchingLines(matcher, input, output)

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

		matched, err := PrintMatchingLines(matcher, input, output)

		assert.NoError(t, err)
		assert.True(t, matched)
		assert.Equal(t, "line1\n\nline2\nline3\n\n", output.String())
	})
}
