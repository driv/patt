package main

import (
	"bufio"
	"os"
	"patt"
)

func main() {
	inputFile := "/tmp/log/syslog"
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	matcher, _ := patt.NewMatcher("<_> pop-os <_>")

	scanner := bufio.NewScanner(file)
	var match bool
	for scanner.Scan() {
		line := scanner.Bytes()
		if matcher.Match(line) {
			match = true
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if !match {
		os.Exit(1)
	}
}
