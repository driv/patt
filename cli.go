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

	var match bool
	if len(params.InputFiles) == 0 {
		processor := NewLineProcessor(io.NopCloser(stdin), stdout, replacer, params.Keep)
		match, err = processor.Process()
		if err != nil {
			return fmt.Errorf("error matching lines: %w", err)
		}
	} else if len(params.InputFiles) == 1 {
		inputFile, err := os.OpenFile(params.InputFiles[0], os.O_RDONLY, 0)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", params.InputFiles[0], err)
		}
		defer inputFile.Close()
		processor := NewLineProcessor(inputFile, stdout, replacer, params.Keep)
		match, err = processor.Process()
		if err != nil {
			return fmt.Errorf("error matching lines: %w", err)
		}
	} else {
		processor := NewFilesProcessor(params.InputFiles, replacer, stdout, params.Keep)
		match, err = processor.Process()
		if err != nil {
			return fmt.Errorf("error matching files: %w", err)
		}
	}
	if !match {
		return fmt.Errorf("no match")
	}

	return nil
}

func replacer(params CLIParams) (LineReplacer, error) {
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
