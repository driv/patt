package patt

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func RunCLI(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	params, err := ParseCLIParams(args[1:])
	if err != nil {
		return fmt.Errorf("bad parameters: %w", err)
	}

	replacer, err := replacer(params)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}

	processor := NewLineProcessor(replacer, params.Keep)

	var match bool
	if len(params.InputFiles) == 0 {
		match, err = processor.Process(ctx, io.NopCloser(stdin), stdout)
		if err != nil {
			return fmt.Errorf("error matching lines: %w", err)
		}
	} else {
		filesProcessor := NewFilesProcessor(params.InputFiles, processor, stdout, &DefaultFileOpener{})
		match, err = filesProcessor.Process(ctx)
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
