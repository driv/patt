package patt

import (
	"fmt"
	"github.com/spf13/cobra"
)

// CLIParams holds the command-line parameters.
type CLIParams struct {
	SearchPatterns  []string
	ReplaceTemplate string
	InputFiles      []string
	Keep            bool
	CPUProfile      string
}

// ParseCLIParams parses flags + positional args
//
//
//	patt [flags] search_pattern [[more_search ...] replace_pattern]
//	     [-- file1 [file2 ...]]
//
// Flags:   -k / --keep  (bool)
func ParseCLIParams(argsWithFlags []string) (CLIParams, error) {
	var out CLIParams

	cmd := &cobra.Command{
		Use:  "patt [flags] search_pattern [[search_pattern ...] replace_template] [-- input_files...]",
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			doubleDashPos := cmd.ArgsLenAtDash()
			var patterns []string
			if doubleDashPos == -1 {
				patterns = args
			} else {
				patterns = args[:doubleDashPos]
				out.InputFiles = args[doubleDashPos:]
			}

			switch len(patterns) {
			case 0:
				return fmt.Errorf("at least one search pattern is required")
			case 1:
				out.SearchPatterns = patterns
			default:
				out.SearchPatterns = patterns[:len(patterns)-1]
				out.ReplaceTemplate = patterns[len(patterns)-1]
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&out.Keep, "keep", "k", false, "print non‑matching lines")
	cmd.Flags().StringVar(&out.CPUProfile, "cpu-profile", "", "write cpu profile to file")
	if err := cmd.Flags().MarkHidden("cpu-profile"); err != nil {
		return out, err
	}

	if err := cmd.ParseFlags(argsWithFlags); err != nil {
		return out, err
	}
	if err := cmd.RunE(cmd, cmd.Flags().Args()); err != nil {
		return out, err
	}
	return out, nil
}
