package patt_test

import (
	"patt"
	"testing"
)

func makeReplacer(t testing.TB, pattern, template string) patt.LineReplacer {
	t.Helper()
	replacer, err := patt.NewReplacer(pattern, template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return replacer
}

func makeMatcher(t testing.TB, stringPattern string) patt.LineReplacer {
	t.Helper()
	matcher, err := patt.NewFilter(stringPattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return matcher
}
