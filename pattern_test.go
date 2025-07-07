package patt_test

import (
	"github.com/google/go-cmp/cmp"
	"patt"
	"testing"
)

func TestMatcher(t *testing.T) {
	tests := []struct {
		name          string
		stringPattern string
		inputLine     string
		match         bool
	}{
		{
			name:          "empty line",
			stringPattern: "<_>",
			inputLine:     "",
			match:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, _ := patt.NewFilter(tt.stringPattern)
			gotMatch := matcher.Match([]byte(tt.inputLine))
			if gotMatch != tt.match {
				t.Errorf("NewMatcher() = %v, expected: %v", gotMatch, tt.match)
			}
		})
	}
}
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
		{
			name:          "Correct escaped ",
			stringPattern: "something \\<_\\> something else",
			wantErr:       false,
			wantNil:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := patt.NewFilter(tt.stringPattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatcher() error = %v, expected err: %v", err, tt.wantErr)
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewMatcher() = %v, expected nil: %v", got, tt.wantNil)
			}
		})
	}
}

func TestReplacer(t *testing.T) {
	tests := []struct {
		name                  string
		stringPattern         string
		stringReplaceTemplate string
		inputLine             string
		expectedResult        string
	}{
		{
			name:                  "Single replacement",
			stringPattern:         "My name is <name>",
			stringReplaceTemplate: "Hello <name>!",
			inputLine:             "My name is Federico",
			expectedResult:        "Hello Federico!",
		},
		{
			name:                  "Double replacement in order",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "Good morning <name> <surname>!",
			expectedResult:        "Good morning Federico Nafria!",
		},
		{
			name:                  "Starts with capture",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "<name> <surname>!",
			expectedResult:        "Federico Nafria!",
		},
		{
			name:                  "Ends with capture",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "<name> <surname>",
			expectedResult:        "Federico Nafria",
		},
		{
			name:                  "Only capture",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "<name>",
			expectedResult:        "Federico",
		},
		{
			name:                  "Double replacement out of order",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "Good morning <surname> <name>!",
			expectedResult:        "Good morning Nafria Federico!",
		},
		{
			name:                  "Duplicated replacement captures",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "Good morning <surname> <name>! May I call you <name>?",
			expectedResult:        "Good morning Nafria Federico! May I call you Federico?",
		},
		{
			name:                  "Consecutive placeholders",
			stringPattern:         "My name is <name> <surname>",
			inputLine:             "My name is Federico Nafria",
			stringReplaceTemplate: "Your username is: <name><surname>",
			expectedResult:        "Your username is: FedericoNafria",
		},
		{
			name:                  "Whitespaces",
			stringPattern:         "    <number> <day>",
			inputLine:             "    284 Mon",
			stringReplaceTemplate: "There were <number> errors on <day>",
			expectedResult:        "There were 284 errors on Mon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacer := makeReplacer(t, tt.stringPattern, tt.stringReplaceTemplate)
			matches := replacer.Match([]byte(tt.inputLine))
			if !matches {
				t.Error("No match")
			}
			replacedString, err := replacer.Replace([]byte(tt.inputLine))
			if err != nil {
				t.Errorf("Error during replacement: %v", err)
			}
			if diff := cmp.Diff(tt.expectedResult, string(replacedString)); diff != "" {
				t.Errorf("Failed Replacement (-expected +got):\n%s", diff)
			}
		})
	}
}

func makeReplacer(t *testing.T, stringPattern string, stringReplaceTemplate string) *patt.Replacer {
	t.Helper()
	replacer, err := patt.NewReplacer(stringPattern, stringReplaceTemplate)
	if err != nil {
		t.Fatalf("Error creating replacer: %v", err)
	}
	return replacer
}

func TestMakeReplacer(t *testing.T) {
	tests := []struct {
		name            string
		matchPattern    string
		replaceTemplate string
		wantErr         bool
	}{
		{
			name:            "Missing placeholder in template",
			matchPattern:    "My name is <name>",
			replaceTemplate: "Hello <name> <surname>!",
			wantErr:         true,
		},
		{
			name:            "Valid template",
			matchPattern:    "My name is <name>",
			replaceTemplate: "Hello <name>!",
		},
		{
			name:            "Invalid match pattern",
			matchPattern:    "My name is <name><surname>",
			replaceTemplate: "Hello <name>!",
			wantErr:         true,
		},
		{
			name:            "Invalid replace template",
			matchPattern:    "My name is <name>",
			replaceTemplate: "Hello <_>!",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := patt.NewReplacer(tt.matchPattern, tt.replaceTemplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewReplacer() error = %v, expected error: %v", err, tt.wantErr)
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Expected an error but got none")
			}
		})
	}
}

func TestMultiReplacer(t *testing.T) {
	tests := []struct {
		name            string
		matchPatterns   []string
		replaceTemplate string
		inputLine       string
		expectedResult  string
		shouldMatch     bool
	}{
		{
			name:            "matches first pattern",
			matchPatterns:   []string{"foo <bar>", "baz <bar>"},
			replaceTemplate: "X:<bar>",
			inputLine:       "foo hello",
			expectedResult:  "X:hello",
			shouldMatch:     true,
		},
		{
			name:            "matches second pattern",
			matchPatterns:   []string{"foo <bar>", "baz <bar>"},
			replaceTemplate: "Y:<bar>",
			inputLine:       "baz world",
			expectedResult:  "Y:world",
			shouldMatch:     true,
		},
		{
			name:            "no match",
			matchPatterns:   []string{"foo <bar>", "baz <bar>"},
			replaceTemplate: "Z:<bar>",
			inputLine:       "no match here",
			expectedResult:  "",
			shouldMatch:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacer := makeMultiReplacer(t, tt.matchPatterns, tt.replaceTemplate)
			matches := replacer.Match([]byte(tt.inputLine))
			if matches != tt.shouldMatch {
				t.Errorf("Match() = %v, want %v", matches, tt.shouldMatch)
			}
			if matches {
				replaced, err := replacer.Replace([]byte(tt.inputLine))
				if err != nil {
					t.Errorf("Replace() error = %v", err)
				}
				if string(replaced) != tt.expectedResult {
					t.Errorf("Replace() = %q, want %q", replaced, tt.expectedResult)
				}
			}
		})
	}
}

func makeMultiReplacer(t *testing.T, patterns []string, template string) *patt.MultiReplacer {
	t.Helper()
	replacer, err := patt.NewMultiReplacer(patterns, template)
	if err != nil {
		t.Fatalf("Error creating MultiReplacer: %v", err)
	}
	return replacer
}

func TestMakeMultiReplacer(t *testing.T) {
	tests := []struct {
		name            string
		matchPatterns   []string
		replaceTemplate string
		wantErr         bool
	}{
		{
			name:            "Valid multi-pattern",
			matchPatterns:   []string{"foo <bar>", "baz <bar>"},
			replaceTemplate: "X:<bar>",
			wantErr:         false,
		},
		{
			name:            "Invalid match pattern",
			matchPatterns:   []string{"foo <bar><baz>"},
			replaceTemplate: "X:<bar>",
			wantErr:         true,
		},
		{
			name:            "Invalid replace template",
			matchPatterns:   []string{"foo <bar>"},
			replaceTemplate: "X:<baz>",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := patt.NewMultiReplacer(tt.matchPatterns, tt.replaceTemplate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMultiReplacer() error = %v, expected error: %v", err, tt.wantErr)
			}
			if tt.wantErr != (err != nil) {
				t.Errorf("Expected an error but got none")
			}
		})
	}
}
