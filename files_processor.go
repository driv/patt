package patt

import (
	"fmt"
	"io"
	"os"
)

// FilesProcessor processes multiple files in parallel and writes the output sequentially.
type FilesProcessor struct {
	files    []string
	replacer LineReplacer
	writer   io.Writer
	keep     bool
}

// NewFilesProcessor creates a new FilesProcessor.
func NewFilesProcessor(files []string, replacer LineReplacer, writer io.Writer, keep bool) *FilesProcessor {
	return &FilesProcessor{
		files:    files,
		replacer: replacer,
		writer:   writer,
		keep:     keep,
	}
}

func (fp *FilesProcessor) Process() (bool, error) {
	result := false
	for _, file := range fp.files {
		inputFile, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			return false, fmt.Errorf("failed to read file %s: %w", file, err)
		}
		defer inputFile.Close()

		processor := NewLineProcessor(inputFile, fp.writer, fp.replacer, fp.keep)
		matched, err := processor.Process()
		if err != nil {
			return false, fmt.Errorf("error processing file %s: %w", file, err)
		}
		result = matched || result
	}
	return result, nil
}
