package patt

import (
	"fmt"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString string
	InputFile     string
}

// ParseCLIParams parses the command-line arguments into CLIParams.
func ParseCLIParams(args []string) (*CLIParams, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}
	params := &CLIParams{
		PatternString: args[1],
	}
	if len(args) > 2 {
		params.InputFile = args[2]
	}
	return params, nil
}
