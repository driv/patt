package patt

import (
	"context"
	"io"
	"iter"
)

// FileOpener defines the interface for opening files.
type FileOpener interface {
	Open(name string) (io.ReadCloser, error)
}

// FilesProcessor processes multiple files.
type FilesProcessor struct {
	files      iter.Seq[string]
	processor  LineProcessor
	writer     io.Writer
	fileOpener FileOpener
	numWorkers int
}

func NewFilesProcessor(files iter.Seq[string], processor LineProcessor, writer io.Writer, fileOpener FileOpener, numWorkers int) *FilesProcessor {
	return &FilesProcessor{
		files:      files,
		processor:  processor,
		writer:     writer,
		fileOpener: fileOpener,
		numWorkers: numWorkers,
	}
}

func (fp *FilesProcessor) Process(ctx context.Context) (result bool, err error) {
	for file := range fp.files {
		rc, err := fp.fileOpener.Open(file)
		if err != nil {
			return false, err
		}
		defer rc.Close()

		matched, err := fp.processor.Process(ctx, rc, fp.writer)
		if err != nil {
			break
		}
		result = result || matched
	}
	return result, nil
}
