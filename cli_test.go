package patt_test

import (
	"bytes"
	"patt"
	"testing"
)

func TestRunCLI(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		stdin     string
		expectErr bool
		expectOut string
	}{
		{
			name:      "match found in stdin",
			args:      []string{"patt", "something <placeholder>", "found <placeholder>!"},
			stdin:     "something match\n",
			expectOut: "found match!\n",
		},
		{
			name:      "no match in stdin",
			args:      []string{"patt", "something <placeholder>", "found <placeholder>!"},
			stdin:     "other match\n",
			expectErr: true,
		},
		{
			name:      "invalid replacement pattern",
			args:      []string{"patt", "something <placeholder>", "found <wrong>!"},
			stdin:     "something match\n",
			expectErr: true,
		},
		{
			name:      "invalid search pattern",
			args:      []string{"patt", "something <placeholder><wrong>", "found <placeholder>!"},
			stdin:     "something match\n",
			expectErr: true,
		},
		{
			name:      "invalid input file",
			args:      []string{"patt", "something <placeholder>", "found <placeholder>!", "testdata/non-existent.log"},
			expectErr: true,
		},
		{
			name:      "no arguments",
			args:      []string{"patt"},
			stdin:     "something to match",
			expectErr: true,
			expectOut: "",
		},
		{
			name:      "missing replacement pattern",
			args:      []string{"patt", "something <placeholder>"},
			stdin:     "something to match",
			expectErr: true,
		},
		{
			name:      "search mode: match found in stdin",
			args:      []string{"patt", "-R", "something <_>"},
			stdin:     "something match\n",
			expectOut: "something match\n",
		},
		{
			name:      "search mode: no match in stdin",
			args:      []string{"patt", "-R", "something <_>"},
			stdin:     "other match\n",
			expectErr: true,
		},
		{
			name:      "search mode: invalid pattern",
			args:      []string{"patt", "-R", "something <_><_>"},
			stdin:     "something match",
			expectErr: true,
		},
		{
			name:      "search mode: missing argument",
			args:      []string{"patt", "-R"},
			stdin:     "something",
			expectErr: true,
		},
		{
			name:      "search mode: match found in file",
			args:      []string{"patt", "-R", "[Sun Dec 04 04:51:08 2005] <_>", "testdata/Apache_2k.log"},
			expectOut: "[Sun Dec 04 04:51:08 2005] [notice] jk2_init() Found child 6725 in scoreboard slot 10\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := bytes.NewReader([]byte(tt.stdin))
			stdout := &bytes.Buffer{}

			err := patt.RunCLI(tt.args, stdin, stdout)

			if (err != nil) != tt.expectErr {
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			}
			if stdout.String() != tt.expectOut {
				t.Errorf("expected stdout %q, got %q", tt.expectOut, stdout.String())
			}
		})
	}
}
