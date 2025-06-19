package patt

import (
	"fmt"
	"io"
	"os"
)

func RunCLI(args []string, stdin io.Reader, stdout io.Writer) error {
	params, err := ParseCLIParams(args[1:])
	if err != nil {
		return fmt.Errorf("bad parameters: %w", err)
	}

	var input io.Reader
	if params.InputFile == "" {
		input = stdin
	} else {
		inputFile, err := os.OpenFile(params.InputFile, os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("error opening input file: %w", err)
		}
		defer inputFile.Close()
		input = inputFile
	}

	var match bool
	if params.ReplaceTemplate == "" {
		filter, err := NewMatcher(params.PatternString)
		if err != nil {
			return fmt.Errorf("cannot parse match template: %w", err)
		}
		match, err = PrintMatchingLines(filter, input, stdout)
		if err != nil {
			return fmt.Errorf("error matching lines: %w", err)
		}
	} else {
		filter, err := NewReplacer(params.PatternString, params.ReplaceTemplate)
		if err != nil {
			return fmt.Errorf("cannot parse template: %w", err)
		}
		match, err = PrintLines(filter, input, stdout)
		if err != nil {
			return fmt.Errorf("error replacing lines: %w", err)
		}
	}
		if !match {
			return fmt.Errorf("no match")
		}

	return nil
}
