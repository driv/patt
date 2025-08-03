package patt

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
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
	startWriting      *sync.Mutex
	doneWriting    *sync.Mutex
	writer  io.Writer
	writing bool
}

func (bw *blockingWriter) Write(p []byte) (n int, err error) {
	if !bw.writing {
		bw.startWriting.Lock()
		bw.writing = true
	}
	return bw.writer.Write(p)
}

type processResult struct {
	matched bool
	err     error
}

func (fp *FilesProcessor) Process(ctx context.Context) (bool, error) {
	resultsChan := make(chan processResult, 10)
	mutexChan := make(chan *blockingWriter, 10)

	go func() {
		for _, file := range fp.files {
			blockedOutput := newBlockedOutput(fp)
			mutexChan <- blockedOutput
			go func(file string) {
				matched, err := fp.processFile(ctx, file, blockedOutput)
				resultsChan <- processResult{
					matched: matched,
					err:     err,
				}
				blockedOutput.doneWriting.Unlock()
			}(file)
		}
		close(mutexChan)
	}()
	finalResultsChan := make(chan processResult, 1)
	go func() {
		var firstErr error
		var anyMatched bool = false
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

	for bo := range mutexChan {
		bo.startWriting.Unlock()
		bo.doneWriting.Lock()
	}
	close(resultsChan)
	finalResult := <-finalResultsChan
	close(finalResultsChan)
	return finalResult.matched, finalResult.err
}

func (fp *FilesProcessor) processFile(ctx context.Context, file string, blockedOutput *blockingWriter) (bool, error) {
	input, err := fp.fileOpener.Open(file)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", file, err)
	}
	defer input.Close()

	matched, err := fp.processor.Process(ctx, input, blockedOutput)
	if err != nil {
		return false, fmt.Errorf("error processing file %s: %w", file, err)
	}
	return matched, nil
}

func newBlockedOutput(fp *FilesProcessor) *blockingWriter {
	blockingOutput := &blockingWriter{writer: fp.writer, startWriting: &sync.Mutex{}, doneWriting: &sync.Mutex{}}
	blockingOutput.startWriting.Lock()
	blockingOutput.doneWriting.Lock()
	return blockingOutput
}
