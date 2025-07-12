package patt

import (
	"reflect"
	"testing"
)

func TestParseCLIParams_NoErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want CLIParams
	}{
		{
			name: "replace with input file",
			args: []string{"pattern", "replacement", "-f", "input.txt"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "replacement",
				InputFile:       "input.txt",
			},
		},
		{
			name: "replace from stdin",
			args: []string{"pattern", "replacement"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "replacement",
			},
		},
		{
			name: "search only from stdin",
			args: []string{"pattern"},
			want: CLIParams{
				SearchPatterns: []string{"pattern"},
			},
		},
		{
			name: "search only with input file",
			args: []string{"pattern", "-f", "input.txt"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "",
				InputFile:       "input.txt",
			},
		},
		{
			name: "replace with input file and keep",
			args: []string{"pattern", "replacement", "-f", "input.txt", "-k"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "replacement",
				InputFile:       "input.txt",
				Keep:            true,
			},
		},
		{
			name: "search only with keep",
			args: []string{"pattern", "-k"},
			want: CLIParams{
				SearchPatterns: []string{"pattern"},
				Keep:           true,
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
			if got != nil && !reflect.DeepEqual(*got, tt.want) {
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
