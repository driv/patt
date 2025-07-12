package patt

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func RunCLI(args []string, stdin io.Reader, stdout io.Writer) error {
	params, err := ParseCLIParams(args[1:])
	if err != nil {
		return fmt.Errorf("bad parameters: %w", err)
	}

	replacer, err := replacer(params)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
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

	processor := NewLineProcessor(input, stdout, params.Keep)

	match, err := processor.ProcessLines(replacer)
	if err != nil {
		return fmt.Errorf("error matching lines: %w", err)
	}
	if !match {
		return fmt.Errorf("no match")
	}

	return nil
}

func replacer(params *CLIParams) (LineReplacer, error) {
	switch {
	case params.ReplaceTemplate == "":
		return NewFilter(params.SearchPatterns[0])
	case len(params.SearchPatterns) == 1:
		return NewReplacer(params.SearchPatterns[0], params.ReplaceTemplate)
	case len(params.SearchPatterns) > 1:
		return NewMultiReplacer(params.SearchPatterns, params.ReplaceTemplate)
	}
	return nil, errors.New("invalid parameters, cannot initialize replacer")
}
