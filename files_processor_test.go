package patt

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"
)

type mockReplacer struct {
	toReplace   string
	replaceWith string
}

func (m *mockReplacer) Replace(line []byte) []byte {
	if strings.Contains(string(line), m.toReplace) {
		return []byte(strings.ReplaceAll(string(line), m.toReplace, m.replaceWith))
	}
	return line
}

func (m *mockReplacer) Match(line []byte) bool {
	return strings.Contains(string(line), m.toReplace)
}

type noopReplacer struct{}

func (n *noopReplacer) Replace(line []byte) []byte {
	return line
}

func (n *noopReplacer) Match(line []byte) bool {
	return true // Always match, but do nothing
}

func TestFilesProcessor_Process_SingleFile(t *testing.T) {
	content := "hello world\n"
	file, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())
	file.WriteString(content)
	file.Close()

	var buf bytes.Buffer
	replacer := &mockReplacer{toReplace: "world", replaceWith: "Go"}
	
	// Use a real LineProcessor
	lineProcessor := NewLineProcessor(replacer, true)

	fp := NewFilesProcessor([]string{file.Name()}, lineProcessor, &buf)

	matched, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !matched {
		t.Errorf("Process() matched = false, want true")
	}
	if !strings.Contains(buf.String(), "hello Go\n") {
		t.Errorf("output = %q, want to contain 'hello Go\n'", buf.String())
	}
}

func TestFilesProcessor_Process_MultipleFiles_Order(t *testing.T) {
	contents := []string{"file1\n", "file2\n", "file3\n"}
	files := make([]string, 3)
	for i, c := range contents {
		f, err := os.CreateTemp("", "testfile-*.txt")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(f.Name())
		f.WriteString(c)
		f.Close()
		files[i] = f.Name()
	}
	var buf bytes.Buffer
	lineProcessor := NewLineProcessor(&noopReplacer{}, true)
	fp := NewFilesProcessor(files, lineProcessor, &buf)

	_, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	out := buf.String()
	for _, c := range contents {
		if !strings.Contains(out, c) {
			t.Errorf("output missing %q", c)
		}
	}
	// Ensure order
	if !strings.Contains(out, "file1\nfile2\nfile3\n") {
		t.Errorf("output order wrong: %q", out)
	}
}