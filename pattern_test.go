package patt_test

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"patt"
	"testing"
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
		{
			name:          "Correct escaped ",
			stringPattern: "something \\<_\\> something else",
			wantErr:       false,
			wantNil:       false,
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

func TestMakeReplacer(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacer, err := patt.NewReplacer(tt.stringPattern, tt.stringReplaceTemplate)
			if err != nil {
				t.Errorf("Error creating replacer: %v", err)
			}
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

func TestMakeReplacer_Fail_Extra_Literal(t *testing.T) {
	stringPattern := "My name is <name>"
	stringReplaceTemplate := "Hello <name> <surname>!"
	_, err := patt.NewReplacer(stringPattern, stringReplaceTemplate)
	if err == nil {
		t.Errorf("Error expected, the template includes <surname>")
	}
	var errorType *patt.ReplaceNameNotFoundError
	if !errors.As(err, &errorType) {
		t.Errorf("Expected error of type ReplaceNameNotFoundError, but got: %v", err)
	}
}
