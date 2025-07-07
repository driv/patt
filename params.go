package patt

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

// CLIParams holds the command-line parameters.

type CLIParams struct {
	SearchPatterns  []string
	ReplaceTemplate string
	InputFiles      string
	HasInputFiles   bool
	Keep            bool
}

func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	app := kingpin.New("patt", "Pattern-based log matcher and replacer")

	posArgs := app.Arg(
		"patterns",
		"Provide one or more search patterns followed by an optional replacement pattern.\n"+
			"Format:\n"+
			"  patt search_pattern\n"+
			"  patt search_pattern1 replace_pattern\n"+
			"  patt search_pattern1 replace_pattern -- file.log\n"+
			"  patt search_pattern1 replace_pattern replace_pattern2 -- file.log\n"+
			"  patt search_pattern1 [[search_pattern2 ... search_patternN] replace_pattern] [-- file.log [file2.log ... filen.log]\n",
	).Required().Strings()

	keep := app.Flag("keep", "Print non-matching lines").Short('k').Bool()
	// inputFile := app.Flag("file", "Input file (optional)").Short('f').String()

	_, err := app.Parse(argsWithFlags)
	if err != nil {
		return nil, err
	}

	result := CLIParams{
		Keep: *keep,
	}
	for _, value := range *posArgs {
		if value == "--" {
			result.HasInputFiles = true
		} else if !result.HasInputFiles {
			result.SearchPatterns = append(result.SearchPatterns, value)
		} else {
			result.InputFiles = value
		}
	}

	if len(result.SearchPatterns) == 0 {
		return nil, fmt.Errorf("at least one search pattern is required")
	}

	if result.HasInputFiles {
		if len(result.InputFiles) == 0 {
			return nil, fmt.Errorf("input files expected after -- ")
		}
	}


	return &result, nil
}
