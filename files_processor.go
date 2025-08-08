package patt

import (
	"context"
	"fmt"
	"io"
	"os"
)

// FileOpener defines the interface for opening files.
type FileOpener interface {
	Open(name string) (io.ReadCloser, error)
}

// DefaultFileOpener is a concrete implementation of FileOpener that uses os.OpenFile.
type DefaultFileOpener struct{}

// Open implements the FileOpener interface.
func (dfo *DefaultFileOpener) Open(name string) (io.ReadCloser, error) {
	return os.OpenFile(name, os.O_RDONLY, 0)
}

// FilesProcessor processes multiple files in parallel and writes the output sequentially.
type FilesProcessor struct {
	files      []string
	processor  LineProcessor
	writer     io.Writer
	fileOpener FileOpener
}

func NewFilesProcessor(files []string, processor LineProcessor, writer io.Writer, fileOpener FileOpener) *FilesProcessor {
	return &FilesProcessor{
		files:      files,
		processor:  processor,
		writer:     writer,
		fileOpener: fileOpener,
	}
}

type blockingWriter struct {
	writer  *io.PipeWriter
	reader  io.Reader
	written bool
	closed  bool
}

func (bw *blockingWriter) Write(p []byte) (n int, err error) {
	if !bw.written {
		bw.written = true
	}
	return bw.writer.Write(p)
}

func (bw *blockingWriter) Close() error {
	bw.closed = true
	return bw.writer.Close()
}

type processResult struct {
	matched bool
	err     error
}

func (fp *FilesProcessor) Process(ctx context.Context) (bool, error) {
	resultsChan := make(chan processResult, 10)
	outputsChan := make(chan *blockingWriter, 10)
	// nextOutput := make(chan *blockingWriter)

	go func() {
		for _, file := range fp.files {
			blockedOutput := newBlockedOutput()
			outputsChan <- blockedOutput
			go func(file string) {
				defer blockedOutput.Close()
				resultsChan <- fp.processFile(ctx, file, blockedOutput)
			}(file)
		}
		close(outputsChan)
	}()
	finalResultsChan := make(chan processResult, 1)
	go func() {
		var firstErr error
		var anyMatched bool
		for result := range resultsChan {
			if firstErr == nil {
				firstErr = result.err
			}
			anyMatched = anyMatched || result.matched
		}
		finalResultsChan <- processResult{
			matched: anyMatched,
			err:     firstErr,
		}
	}()

	for bo := range outputsChan {
		if !bo.written && bo.closed {
			continue // skip if nothing was written and the writer is closed
		}
		io.Copy(fp.writer, bo.reader)
	}
	close(resultsChan)
	finalResult := <-finalResultsChan
	close(finalResultsChan)
	return finalResult.matched, finalResult.err
}

func (fp *FilesProcessor) processFile(ctx context.Context, file string, blockedOutput *blockingWriter) (result processResult) {
	input, err := fp.fileOpener.Open(file)
	if err != nil {
		return processResult{
			err: fmt.Errorf("failed to read file %s: %w", file, err),
		}
	}
	defer input.Close()

	matched, err := fp.processor.Process(ctx, input, blockedOutput)
	if err != nil {
		return processResult{
			matched: matched,
			err:     fmt.Errorf("error processing file %s: %w", file, err),
		}
	}
	return processResult{matched: matched}
}

func newBlockedOutput() *blockingWriter {
	pr, pw := io.Pipe()
	blockingOutput := &blockingWriter{
		writer: pw,
		reader: pr,
	}
	return blockingOutput
}
