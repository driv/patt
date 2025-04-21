package patt

import (
	"fmt"
	"io"
	"os"
)

func RunCLI(patternString, inputFile, outputFile string) error {
	var input *os.File
	var err error

	if inputFile == "" {
		input = os.Stdin
	} else {
		input, err = os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("error opening input file: %v", err)
		}
		defer input.Close()
	}

	var output *os.File
	if outputFile == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer output.Close()
	}

	return RunCLIWithIO(patternString, input, output)
}

func RunCLIWithIO(patternString string, input io.Reader, output io.Writer) error {
	filter, err := NewMatcher(patternString)
	if err != nil {
		return fmt.Errorf("error creating matcher: %v", err)
	}

	match, err := MatchLines(filter, input, output)
	if err != nil {
		return fmt.Errorf("error matching lines: %v", err)
	}
	if !match {
		return fmt.Errorf("no match")
	}
	return nil
}
