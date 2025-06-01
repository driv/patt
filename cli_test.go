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
			name:      "valid pattern with match",
			args:      []string{"patt", "-R", "something <_>"},
			stdin:     "something match\n",
			expectErr: false,
			expectOut: "something match\n",
		},
		{
			name:      "valid pattern with no match",
			args:      []string{"patt", "-R", "something <_>"},
			stdin:     "other match\n",
			expectErr:	true,
			expectOut: "",
		},
		{
			name:      "invalid pattern",
			args:      []string{"patt", "<_><_>"},
			stdin:     "",
			expectErr: true,
			expectOut: "",
		},
		{
			name:      "missing arg",
			args:      []string{"patt"},
			stdin:     "",
			expectErr: true,
			expectOut: "",
		},
		{
			name:      "search only, input file",
			args:      []string{"patt",  "-R", "[Sun Dec 04 04:51:08 2005] <_>","test_files/Apache_2k.log"},
			stdin:     "",
			expectErr: false,
			expectOut: `[Sun Dec 04 04:51:08 2005] [notice] jk2_init() Found child 6725 in scoreboard slot 10`+"\n",
		},
		{
			name:      "invalid input file",
			args:      []string{"patt", "<_>", "", "test_files/non-existent.log"},
			stdin:     "",
			expectErr: true,
			expectOut: "",
		},
		{
			name:      "valid replacement with match",
			args:      []string{"patt", "something <placeholder>", "found <placeholder>!"},
			stdin:     "something match\n",
			expectErr: false,
			expectOut: "found match!\n",
		},
		{
			name:      "valid replacement with no match",
			args:      []string{"patt", "something <placeholder>", "found <placeholder>!"},
			stdin:     "other match\n",
			expectErr:	true,
			expectOut: "",
		},
		{
			name:      "invalid replacement",
			args:      []string{"patt", "something <placeholder>", "found <wrong>!"},
			stdin:     "",
			expectErr: true,
			expectOut: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := bytes.NewReader([]byte(tt.stdin))
			stdout := &bytes.Buffer{}

			err := patt.RunCLI(tt.args, stdin, stdout)

			if (err != nil) != tt.expectErr{
				t.Errorf("expected error %v, got %v", tt.expectErr, err)
			} 
			if stdout.String() != tt.expectOut {
				t.Errorf("expected stdout %q, got %q", tt.expectOut, stdout.String())
			}
		})
	}
}
