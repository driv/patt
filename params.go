package patt

import (
	"errors"
	"flag"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString     string
	InputFile         string
	ReplacementString string
}

// Updated ParseCLIParams to handle ReplacementString
func ParseCLIParams(args []string) (*CLIParams, error) {
	flags := flag.NewFlagSet("patt", flag.ContinueOnError)
	flags.Parse(args)

	params := &CLIParams{
		PatternString:     flags.Arg(0),
		ReplacementString: flags.Arg(1),
		InputFile:         flags.Arg(2),
	}

	if params.PatternString == "" {
		return nil, errors.New("patt match_pattern replace_pattern [file]")
	}

	return params, nil
}
