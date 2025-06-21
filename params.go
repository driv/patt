package patt

import (
	"github.com/alecthomas/kingpin/v2"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString   string
	ReplaceTemplate string
	InputFile       string
	Keep            bool
}

func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	app := kingpin.New("patt", "Pattern-based log matcher and replacer")
	pattern := app.Arg("pattern", "Pattern to search for").Required().String()
	replacement := app.Arg("replacement", "Replacement template (optional)").Default("").String()
	inputFile := app.Flag("file", "Input file (optional)").Short('f').String()
	keep := app.Flag("keep", "Print non-matching lines").Short('k').Bool()

	_, err := app.Parse(argsWithFlags)
	if err != nil {
		return nil, err
	}

	return &CLIParams{
		PatternString:   *pattern,
		ReplaceTemplate: *replacement,
		InputFile:       *inputFile,
		Keep:            *keep,
	}, nil
}
