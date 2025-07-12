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
			args: []string{"pattern", "replacement", "--", "input.txt"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "replacement",
				InputFiles:      []string{"input.txt"},
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
			args: []string{"pattern", "--", "input.txt"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "",
				InputFiles:      []string{"input.txt"},
			},
		},
		{
			name: "search only with 2 input files",
			args: []string{"pattern", "--", "input.txt", "input2.txt"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "",
				InputFiles:      []string{"input.txt", "input2.txt"},
			},
		},
		{
			name: "replace with input file and keep",
			args: []string{"pattern", "replacement", "--", "input.txt", "-k"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern"},
				ReplaceTemplate: "replacement",
				InputFiles:      []string{"input.txt"},
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
		{
			name: "multiple search patterns",
			args: []string{"pattern1", "pattern2", "template"},
			want: CLIParams{
				SearchPatterns:  []string{"pattern1", "pattern2"},
				ReplaceTemplate: "template",
				Keep:            false,
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
