package patt

import "errors"

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString     string
	InputFile         string
	ReplacementString string
}

// Updated ParseCLIParams to handle ReplacementString
func ParseCLIParams(args []string) (*CLIParams, error) {
	v := func(a []string, i int) string {
		if !(len(a) > i) {
			return ""
		}
		return a[i]
	}
	result := &CLIParams{
		PatternString:     v(args, 0),
		ReplacementString: v(args, 1),
		InputFile:         v(args, 2),
	}

	if result.PatternString == "" {
		return nil, errors.New("patt match_pattern replace_pattern [file]")
	}
	return result, nil
}
