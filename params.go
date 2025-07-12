package patt

import (
	"fmt"

	"github.com/alecthomas/kingpin/v2"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	SearchPatterns  []string
	ReplaceTemplate string
	InputFile       string
	Keep            bool
}

type patterns []string

func (p *patterns) String() string {
	return ""
}

func (p *patterns) Set(value string) error {
	*p = append(*p, value)
	return nil
}

func (p *patterns) IsCumulative() bool {
	return true
}

func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	app := kingpin.New("patt", "Pattern-based log matcher and replacer")

	keep := app.Flag("keep", "Print non-matching lines").Short('k').Bool()
	inputFile := app.Flag("file", "Input file (optional)").Short('f').String()
	patterns := &patterns{}
	app.Arg(
		"patterns",
		"Provide one or more search patterns followed by an optional replacement pattern.\n"+
			"Format:\n"+
			"  patt search_pattern\n"+
			"  patt search_pattern1 replace_pattern\n"+
			"  patt search_pattern1 replace_pattern -f file.log\n"+
			"  patt search_pattern1 replace_pattern replace_pattern2 -- file.log\n"+
			"  patt search_pattern1 [[search_pattern2 ... search_patternN] replace_pattern] [-- file.log [file2.log ... filen.log]\n",
	).SetValue(patterns)

	_, err := app.Parse(argsWithFlags)
	if err != nil {
		return nil, err
	}
	posArgs := *patterns

	result := CLIParams{
		Keep:      *keep,
		InputFile: *inputFile,
	}

	if len(posArgs) == 0 {
		return nil, fmt.Errorf("at least one search pattern is required")
	} else if len(posArgs) == 1 {
		result.SearchPatterns = posArgs
	} else if len(posArgs) > 1 {
		result.SearchPatterns = posArgs[:1]
		result.ReplaceTemplate = posArgs[len(posArgs)-1]
	}

	return &result, nil
}
