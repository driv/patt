package patt

import (
	"errors"
	"flag"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString     string
	ReplacementString string
	InputFile         string
	SearchOnly        bool
}

// Updated ParseCLIParams to handle ReplacementString
func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	flags := flag.NewFlagSet("patt", flag.ContinueOnError)
	var searchOnly bool
	flags.BoolVar(&searchOnly, "R", false, "Search only, without replacement")
	err := flags.Parse(argsWithFlags)
	if err != nil {
		return nil, err
	}

	result := &CLIParams{
		PatternString: flags.Arg(0),
		SearchOnly:    searchOnly,
	}

	if result.PatternString == "" {
		return nil, errors.New("patt match_pattern replace_pattern [file]")
	}

	if searchOnly {
		result.InputFile = flags.Arg(1)
	} else {
		result.ReplacementString = flags.Arg(1)
		result.InputFile = flags.Arg(2)
		if result.ReplacementString == "" {
			return nil, errors.New("replacement pattern not provided")
		}
	}
	return result, nil
}
