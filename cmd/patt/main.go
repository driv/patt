package main

import (
	"os"
	"patt"
)

func main() {
	exitIfErr(patt.RunCLI(os.Args, os.Stdin, os.Stdout))
}

func exitIfErr(err error) {
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
