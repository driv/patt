package patt

import (
	"fmt"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	SearchPatterns  []string
	ReplaceTemplate string
	InputFiles      []string
	Keep            bool
}

func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	result := CLIParams{}
	var positionalArgs []string

	for _, arg := range argsWithFlags {
		if arg == "-k" || arg == "--keep" {
			result.Keep = true
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	var patterns []string
	var files []string
	var hasInputFiles bool
	for _, arg := range positionalArgs {
		if arg == "--" {
			hasInputFiles = true
			continue
		}
		if !hasInputFiles {
			patterns = append(patterns, arg)
		} else {
			files = append(files, arg)
		}
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("at least one search pattern is required")
	}

	if len(patterns) > 1 {
		result.ReplaceTemplate = patterns[len(patterns)-1]
		result.SearchPatterns = patterns[:len(patterns)-1]
	} else {
		result.SearchPatterns = patterns
	}

	if len(files) > 0 {
		result.InputFiles = files
	}

	return &result, nil
}
