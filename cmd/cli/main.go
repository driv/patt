package main

import (
	"fmt"
	"os"
	"patt"
)

func main() {
	params, err := patt.ParseCLIParams(os.Args)
	if err != nil {
		exitIfErr(err)
	}
	filter, err := patt.NewMatcher(params.PatternString)
	if err != nil {
		exitIfErr(fmt.Errorf("error creating matcher: %v", err))
	}

	var input *os.File
	if params.InputFile == "" {
		input = os.Stdin
	} else {
		input, err = os.OpenFile(params.InputFile, os.O_RDONLY, 0)
		if err != nil {
			exitIfErr(fmt.Errorf("error opening input file: %v", err))
		}
		defer input.Close()
	}

	match, err := patt.PrintMatchingLines(filter, input, os.Stdout)
	if err != nil {
		exitIfErr(fmt.Errorf("error matching lines: %v", err))
	}

	if !match {
		exitIfErr(fmt.Errorf("no match"))
	}
}

func exitIfErr(err error) {
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
