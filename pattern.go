package patt

import (
	"github.com/grafana/loki/v3/pkg/logql/log/pattern"
)

type LinesMatcher interface {
	Match(b []byte) bool
}

type PatternMatcher struct {
	filter pattern.Matcher
}

func (m PatternMatcher) Match(b []byte) bool {
	return m.filter.Test(b)
}

func NewMatcher(stringPattern string) (LinesMatcher, error) {
	filter, err := pattern.ParseLineFilter([]byte(stringPattern))
	if err != nil {
		return nil, err
	}
	return PatternMatcher{filter: *filter}, err
}
