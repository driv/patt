package patt

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
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

func newMockFileOpener(files map[string]string) *mockFileOpener {
	mockFiles := make(map[string]io.ReadCloser)
	for name, content := range files {
		mockFiles[name] = io.NopCloser(strings.NewReader(content))
	}
	return &mockFileOpener{
		files: mockFiles,
	}
}

func TestFilesProcessor_Process_SingleFile(t *testing.T) {
	fileName := "testfile.txt"
	fileContent := "hello world\n"

	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			w.Write([]byte("hello Go\n"))
			return true, nil
		},
	}

	fileOpener := newMockFileOpener(map[string]string{
		fileName: fileContent,
	})

	fp := NewFilesProcessor([]string{fileName}, lineProcessor, &buf, fileOpener)

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
	fileContents := map[string]string{
		"file1.txt": "content file 1\n",
		"file2.txt": "content file 2\n",
		"file3.txt": "content file 3\n",
	}
	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}

	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			_, err := io.Copy(w, r)
			return true, err
		},
	}

	fileOpener := newMockFileOpener(fileContents)

	fp := NewFilesProcessor(fileNames, lineProcessor, &buf, fileOpener)

	_, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	out := buf.String()
	expected := "content file 1\ncontent file 2\ncontent file 3\n"
	if out != expected {
		t.Errorf("output = %q, want %q", out, expected)
	}
}

func TestFilesProcessor_Process_MultipleFiles_SingleMatch(t *testing.T) {
	fileContents := map[string]string{
		"file1.txt": "content file 1\n",
		"file2.txt": "content file 2\n",
		"file3.txt": "content file 3\n",
	}
	fileNames := []string{"file1.txt", "file2.txt", "file3.txt"}

	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			content, _ := io.ReadAll(r)
			if string(content) == fileContents["file2.txt"] {
				w.Write(content)
				return true, nil
			}
			return false, nil
		},
	}

	fileOpener := newMockFileOpener(fileContents)

	fp := NewFilesProcessor(fileNames, lineProcessor, &buf, fileOpener)

	_, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	out := buf.String()
	expected := "content file 2\n"
	if out != expected {
		t.Errorf("output = %q, want %q", out, expected)
	}
}

func TestFilesProcessor_Process_FileNotFound(t *testing.T) {
	fileName := "nonexistent.txt"

	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			return false, nil
		},
	}

	fileOpener := newMockFileOpener(map[string]string{})

	fp := NewFilesProcessor([]string{fileName}, lineProcessor, &buf, fileOpener)

	result, err := fp.Process(context.Background())
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
	if result {
		t.Errorf("expected result to be false, got true")
	}
}

func TestFilesProcessor_Process_Files_Written_In_Order(t *testing.T) {
	fileContents := map[string]string{}
	for i := range 5 { //TODO we should fill the buffers
		fileContents[fmt.Sprintf("file%d.txt", i+1)] = fmt.Sprintf("content file %d\n", i+1)
	}
	fileNames := getKeys(fileContents)

	lockFile1 := make(chan struct{})
	// lockFile2 := make(chan struct{})
	doneReading := make(chan struct{})
	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			content, _ := io.ReadAll(r)
			if string(content) == fileContents[fileNames[1]] {
				lockFile1 <- struct{}{} // simulating long processing time
			}
			doneReading <- struct{}{}
			w.Write(content)
			return true, nil
		},
	}

	fileOpener := newMockFileOpener(fileContents)

	fp := NewFilesProcessor(fileNames, lineProcessor, &buf, fileOpener)

	go func() {
		for i := range fileNames {
			if i < 1 {
				continue // skip one file being processed
			}
			<-doneReading // wait for the rest of the files to be processed
		}

		// <-lockFile2  // finish processing first file
		<-lockFile1  // finish processing second file
		for range doneReading {
			<-doneReading // drain the channel to ensure all processing is done
		}
	}()

	result, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !result {
		t.Errorf("expected result to be true, got false")
	}
	out := buf.String()
	var expected strings.Builder
	for _, fileName := range fileNames {
		expected.WriteString(fileContents[fileName])
	}
	if out != expected.String() {
		t.Errorf("output = %q, want %q", out, expected.String())
	}
}

func TestFilesProcessor_Process_Files_Read_In_Parallel(t *testing.T) {
	fileContents := map[string]string{}
	for i := range 5 { //TODO we should fill the buffers
		fileContents[fmt.Sprintf("file%d.txt", i+1)] = fmt.Sprintf("content file %d\n", i+1)
	}
	fileNames := getKeys(fileContents)

	lockFile1 := make(chan struct{})
	doneReading := make(chan struct{})
	var buf bytes.Buffer
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			content, _ := io.ReadAll(r)
			if string(content) == fileContents[fileNames[1]] {
				lockFile1 <- struct{}{} // simulating long processing time
			}
			doneReading <- struct{}{}
			return false, nil
		},
	}

	fileOpener := newMockFileOpener(fileContents)

	fp := NewFilesProcessor(fileNames, lineProcessor, &buf, fileOpener)

	// We should be able to read all files in parallel, without having to wait for the first one to finish.
	go func() {
		for i := range fileNames {
			if i == 0 {
				continue // skip one file that is being processed
			}
			<-doneReading // wait for the rest of the files to be processed
		}

		<-lockFile1 // finish processing file1.txt
		<-doneReading
	}()

	result, err := fp.Process(context.Background())
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if result {
		t.Errorf("expected result to be false, got true")
	}
	out := buf.String()
	if out != "" {
		t.Errorf("nothing should have been written to output, got %q", out)
	}

}

func getKeys(fileContents map[string]string) []string {
	fileNames := make([]string, len(fileContents))
	i := 0
	for k := range fileContents {
		fileNames[i] = k
		i++
	}
	return fileNames
}

//This should be a benchmark test.
func BenchmarkFilesProcessor_Process(b *testing.B) {
	fileContents := map[string]string{}
	for i := range 1_000 {
		fileContents[fmt.Sprintf("file%d.txt", i+1)] = fmt.Sprintf("content file %d\n", i+1)
	}
	fileNames := getKeys(fileContents)

	lockFile1 := make(chan struct{})
	lockFile2 := make(chan struct{})
	processing := make(chan struct{})
	lineProcessor := &mockLineProcessor{
		processFunc: func(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
			content, _ := io.ReadAll(r)
			if string(content) == fileContents[fileNames[1]] {
				lockFile1 <- struct{}{} // simulating long processing time
			}
			if string(content) == fileContents[fileNames[2]] {
				lockFile2 <- struct{}{} // simulating long processing time
			}
			processing <- struct{}{}
			// w.Write(content)
			//write half the content to simulate partial processing
			w.Write([]byte(string(content)[:len(content)/2]))
			// sleep for a few milliseconds to simulate processing time
			time.Sleep(10 * time.Millisecond)
			// write the rest of the content
			w.Write([]byte(string(content)[len(content)/2:]))
			return true, nil
		},
	}

	fileOpener := newMockFileOpener(fileContents)

	fp := NewFilesProcessor(fileNames, lineProcessor, io.Discard, fileOpener)

	go func() {
		<-processing // wait for some files to be processed
		<-processing // wait for some files to be processed
		<-processing // wait for some files to be processed
		<-lockFile2  // finish processing first file
		<-lockFile1  // finish processing second file
		for range processing {
			<-processing // drain the channel to ensure all processing is done
		}
	}()

	result, err := fp.Process(context.Background())
	if err != nil {
		b.Fatalf("Process() error = %v", err)
	}
	if !result {
		b.Errorf("expected result to be true, got false")
	}
}