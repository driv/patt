package main

import (
	"os"
	patt "pattern"

)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	patternString := os.Args[1]
	inputFile := ""
	if len(os.Args) > 2 {
		inputFile = os.Args[2]
	}
	err := patt.RunCLI(patternString, inputFile, "")
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

