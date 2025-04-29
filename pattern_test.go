package patt_test

import (
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
