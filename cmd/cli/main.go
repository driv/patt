package main

import (
	"fmt"
	"io"
	"log"
	"os"
	patt "pattern"
)

func main() {
	pattern, inputFile := parseArgs(os.Args)
	input, err := openInputFile(inputFile)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer input.Close()

	err = processPattern(pattern, input)
	if err != nil {
		log.Fatalf("error processing pattern: %v", err)
	}
}

func openInputFile(inputFile string) (*os.File, error) {
	if inputFile == "" {
		return os.Stdin, nil
	}
	return os.Open(inputFile)
}

func parseArgs(args []string) (string, string) {
	if len(args) < 2 {
		log.Fatalf("usage: %s <pattern> [input-file]", args[0])
	}
	pattern := args[1]
	inputFile := ""
	if len(args) > 2 {
		inputFile = args[2]
	}
	return pattern, inputFile
}

func processPattern(pattern string, input io.Reader) error {
	filter, err := patt.NewMatcher(pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}
	matched, err := patt.MatchLines(filter, input, os.Stdout)
	if err != nil {
		return err
	}
	if !matched {
		_, _ = fmt.Fprintln(os.Stderr, "no matches found")
		os.Exit(1)
	}
	return nil
}
