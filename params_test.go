package patt

import (
	"testing"
)

func TestParseCLIParams_NoErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want *CLIParams
	}{
		{
			name: "valid arguments",
			args: []string{"program", "pattern", "input.txt"},
			want: &CLIParams{PatternString: "pattern", InputFile: "input.txt"},
		},
		{
			name: "missing input file",
			args: []string{"program", "pattern"},
			want: &CLIParams{PatternString: "pattern"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCLIParams(tt.args)
			if err != nil {
				t.Errorf("ParseCLIParams() error = %v, want no error", err)
				return
			}
			if got != nil && *got != *tt.want {
				t.Errorf("ParseCLIParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCLIParams_WithErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "not enough arguments",
			args:    []string{"program"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCLIParams(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCLIParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
