package patt

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"slices"
)

func RunCLI(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer) error {
	params, err := ParseCLIParams(args[1:])
	if err != nil {
		return fmt.Errorf("bad parameters: %w", err)
	}

	if params.CPUProfile != "" {
		f, err := os.Create(params.CPUProfile)
		if err != nil {
			return fmt.Errorf("could not create CPU profile: %w", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			return fmt.Errorf("could not start CPU profile: %w", err)
		}
		defer pprof.StopCPUProfile()
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
	} else if len(params.InputFiles) == 1 {
		fileOpener := &BufferedFileOpener{}
		rc, err := fileOpener.Open(params.InputFiles[0])
		if err != nil {
			return fmt.Errorf("cannot open input file: %w", err)
		}
		defer rc.Close()

		match, err = processor.Process(ctx, rc, stdout)
		if err != nil {
			return fmt.Errorf("error matching file: %w", err)
		}
	} else {
		filesProcessor := NewFilesProcessor(
			slices.Values(params.InputFiles),
			processor,
			stdout,
			&BufferedFileOpener{},
			4,
		)
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

type BufferedFileOpener struct {
	BufSize int
}

func (ffo *BufferedFileOpener) Open(name string) (io.ReadCloser, error) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	if ffo.BufSize <= 0 {
		ffo.BufSize = 4 * 1024 * 1024 // default 4 MB
	}
	br := bufio.NewReaderSize(f, ffo.BufSize)

	return struct {
		io.Reader
		io.Closer
	}{
		Reader: br,
		Closer: f,
	}, nil
}
