package patt

import (
	"fmt"

	"patt/pattern"
)

type ReplaceNameNotFoundError struct {
	Name string
}

func (e *ReplaceNameNotFoundError) Error() string {
	return fmt.Sprintf("replace name '%s' not found in source names", e.Name)
}

type LinesMatcher interface {
	Match(b []byte) bool
}

type PatternMatcher struct {
	filter pattern.Matcher
}

func (m PatternMatcher) Match(b []byte) bool {
	return m.filter.Test(b)
}

func NewFilter(stringPattern string) (LineReplacer, error) {
	filter, err := pattern.ParseLineFilter([]byte(stringPattern))
	if err != nil {
		return nil, err
	}
	matcher := PatternMatcher{filter: *filter}
	replacer := matchFilter{PatternMatcher: &matcher}
	return replacer, nil
}

type matchFilter struct {
	*PatternMatcher
}

func (mf matchFilter) Replace(line []byte) ([]byte, error) {
	return line, nil
}

func NewReplacer(stringPattern, stringReplaceTemplate string) (*Replacer, error) {
	filter, err := pattern.New(stringPattern)
	if err != nil {
		return nil, err
	}
	sourceCaptures := filter.Names()
	literals, replaceCaptures, err := pattern.ParseNodes(stringReplaceTemplate)
	if err != nil {
		return nil, err
	}
	positions, err := capturesPositions(sourceCaptures, replaceCaptures, literals)
	if err != nil {
		return nil, err
	}
	return &Replacer{
		PatternMatcher: &PatternMatcher{filter: *filter},
		literals:       literals,
		positions:      positions,
	}, nil
}

func capturesPositions(sourceNames []string, replaceNames []string, literals [][]byte) ([]int, error) {
	sourceNameSet := make(map[string]int)
	for pos, name := range sourceNames {
		sourceNameSet[name] = pos
	}
	positions := make([]int, len(replaceNames))
	for i, replaceName := range replaceNames {
		if literals[i] != nil {
			continue
		}
		pos, exists := sourceNameSet[replaceName]
		if !exists {
			return nil, &ReplaceNameNotFoundError{Name: replaceName}
		}
		positions[i] = pos
	}
	return positions, nil
}

type LineReplacer interface {
	LinesMatcher
	Replace(b []byte) ([]byte, error)
}
type Replacer struct {
	*PatternMatcher
	literals  [][]byte
	positions []int
}

func (r *Replacer) Replace(b []byte) ([]byte, error) {
	matches := r.filter.Matches(b)
	var result []byte
	for i, l := range r.literals {
		if l != nil {
			result = append(result, l...)
		} else {
			result = append(result, matches[r.positions[i]]...)
		}
	}

	return result, nil
}
