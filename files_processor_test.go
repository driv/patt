package patt

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
)

// mockLineProcessor is a mock implementation of the LineProcessor interface.
type mockLineProcessor struct {
	// processFunc allows for custom logic in the Process method for different test cases.
	processFunc func(ctx context.Context, r io.Reader, w io.Writer) (bool, error)
}

// Process delegates the call to the processFunc.
func (m *mockLineProcessor) Process(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
	if m.processFunc != nil {
		return m.processFunc(ctx, r, w)
	}
	// Default behavior: do nothing.
	return false, nil
}

// mockFileOpener is a mock implementation of the FileOpener interface for testing.
type mockFileOpener struct {
	files map[string]io.ReadCloser
}

// Open implements the FileOpener interface, returning a mock file content.
func (m *mockFileOpener) Open(name string) (io.ReadCloser, error) {
	if content, ok := m.files[name]; ok {
		return content, nil
	}
	return nil, os.ErrNotExist
}

func memoryFilesOpener(files map[string]string) *mockFileOpener {
	mockFiles := make(map[string]io.ReadCloser)
	for name, content := range files {
		mockFiles[name] = io.NopCloser(strings.NewReader(content))
	}
	return &mockFileOpener{
		files: mockFiles,
	}
}

const numWorkers = 3

func TestFilesProcessor_Process_SingleFile_ReadOnly(t *testing.T) {
	fileNames, files := makeFiles(1)

	processFunc := func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
		return false, nil
	}

	fp := NewFilesProcessor(
		slices.Values(fileNames),
		makeMockLineProcessor(processFunc),
		nil,
		memoryFilesOpener(files),
		numWorkers,
	)

	matched, err := fp.Process(t.Context())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if matched {
		t.Errorf("Process() matched = %v, want false", matched)
	}
}
func TestFilesProcessor_Process_SingleFile(t *testing.T) {
	fileNames, files := makeFiles(1)

	processFunc := func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
		_, err := io.Copy(w, r)
		return true, err
	}
	var buffer bytes.Buffer

	fp := NewFilesProcessor(
		slices.Values(fileNames),
		makeMockLineProcessor(processFunc),
		&buffer,
		memoryFilesOpener(files),
		numWorkers,
	)

	matched, err := fp.Process(t.Context())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !matched {
		t.Errorf("Process() matched = false, want true")
	}
	if buffer.String() != files[fileNames[0]] {
		t.Errorf("output = %q, expected '%s'", buffer.String(), files[fileNames[0]])
	}
}

func makeFiles(n int) ([]string, map[string]string) {
	fileNames := make([]string, n)
	files := make(map[string]string, n)
	for i := range n {
		fileName := fmt.Sprintf("testfile%d.txt", i)
		fileNames[i] = fileName
		files[fileName] = fmt.Sprintf("hello world %d\n", i)
	}
	return fileNames, files
}

func makeMockLineProcessor(processFunc func(ctx context.Context, r io.Reader, w io.Writer) (bool, error)) *mockLineProcessor {
	lineProcessor := &mockLineProcessor{
		processFunc: processFunc,
	}
	return lineProcessor
}

func TestFilesProcessor_Process_MultipleFiles_Order(t *testing.T) {
	fileNames, fileContents := makeFiles(3)
	fileOpener := memoryFilesOpener(fileContents)

	var buf bytes.Buffer
	processFunc := func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
		_, err := io.Copy(w, r)
		return true, err
	}

	fp := NewFilesProcessor(
		slices.Values(fileNames),
		makeMockLineProcessor(processFunc),
		&buf,
		fileOpener,
		numWorkers,
	)

	matched, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !matched {
		t.Errorf("Process() matched = false, want true")
	}
	var expected strings.Builder
	for _, f := range fileNames {
		expected.WriteString(fileContents[f])
	}
	if buf.String() != expected.String() {
		t.Errorf("output = %q, expected '%s'", buf.String(), expected.String())
	}
}

func TestFilesProcessor_Process_MultipleFiles_SingleMatch(t *testing.T) {
	fileNames, fileContents := makeFiles(3)
	fileOpener := memoryFilesOpener(fileContents)

	var buf bytes.Buffer
	processFunc := func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
		content, _ := io.ReadAll(r)
		if string(content) == fileContents[fileNames[2]] {
			w.Write(content)
			return true, nil
		}
		return false, nil
	}

	fp := NewFilesProcessor(
		slices.Values(fileNames),
		makeMockLineProcessor(processFunc),
		&buf,
		fileOpener,
		numWorkers,
	)

	matched, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !matched {
		t.Errorf("Process() matched = false, want true")
	}
	out := buf.String()
	if out != fileContents[fileNames[2]] {
		t.Errorf("output = %q, want %q", out, fileContents[fileNames[2]])
	}
}

func TestFilesProcessor_Process_FileNotFound(t *testing.T) {
	fileName := "nonexistent.txt"
	fileOpener := memoryFilesOpener(map[string]string{})

	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			return false, nil
		},
	}

	fp := NewFilesProcessor(
		slices.Values([]string{fileName}),
		lineProcessor,
		&buf,
		fileOpener,
		numWorkers,
	)

	result, err := fp.Process(context.Background())
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if result {
		t.Errorf("expected result to be false, got true")
	}
	if buf.String() != "" {
		t.Errorf("expected buffer to be empty, got %q", buf.String())
	}
}
