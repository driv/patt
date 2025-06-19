package patt

import (
	"testing"
)

func TestParseCLIParams_NoErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want CLIParams
	}{
		{
			name: "input file",
			args: []string{"pattern", "replacement", "input.txt"},
			want: CLIParams{
				PatternString:   "pattern",
				ReplaceTemplate: "replacement",
				SearchOnly:      false,
				InputFile:       "input.txt",
			},
		},
		{
			name: "stdin",
			args: []string{"pattern", "replacement"},
			want: CLIParams{
				PatternString:   "pattern",
				ReplaceTemplate: "replacement",
				SearchOnly:      false,
				InputFile:       "",
			},
		},
		{
			name: "search only stdin",
			args: []string{"-R", "pattern"},
			want: CLIParams{
				PatternString:   "pattern",
				SearchOnly:      true,
				ReplaceTemplate: "",
				InputFile:       "",
			},
		},
		{
			name: "search only input file",
			args: []string{"-R", "pattern", "input.txt"},
			want: CLIParams{
				PatternString:   "pattern",
				ReplaceTemplate: "",
				SearchOnly:      true,
				InputFile:       "input.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCLIParams(tt.args)
			if err != nil {
				t.Errorf("ParseCLIParams() error = %v, want no error", err)
				return
			}
			if got != nil && *got != tt.want {
				t.Errorf("ParseCLIParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCLIParams_WithErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing pattern",
			args: []string{},
		},
		{
			name: "missing replace template",
			args: []string{"pattern"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCLIParams(tt.args)
			if err == nil {
				t.Errorf("ParseCLIParams(%v) should fail", tt.args)
			}
		})
	}
}
