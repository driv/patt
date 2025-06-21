package patt

import (
	"github.com/alecthomas/kingpin/v2"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	PatternString   string
	ReplaceTemplate string
	InputFile       string
}

// ParseCLIParams using kingpin for robust CLI parsing
func ParseCLIParams(argsWithFlags []string) (*CLIParams, error) {
	app := kingpin.New("patt", "Pattern-based log matcher and replacer")
	pattern := app.Arg("pattern", "Pattern to search for").Required().String()
	replacement := app.Arg("replacement", "Replacement template (optional)").Default("").String()
	inputFile := app.Flag("file", "Input file (optional)").Short('f').String()

	// kingpin expects os.Args[1:], but we want to support custom args
	_, err := app.Parse(argsWithFlags)
	if err != nil {
		return nil, err
	}

	params := &CLIParams{
		PatternString:   *pattern,
		ReplaceTemplate: *replacement,
		InputFile:       *inputFile,
	}
	return params, nil
}
