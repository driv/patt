package patt

import (
	"fmt"
	"io"
	"os"
)

// FilesProcessor processes multiple files in parallel and writes the output sequentially.
type FilesProcessor struct {
	files     []string
	processor *LineProcessor
	writer    io.Writer
}

// NewFilesProcessor creates a new FilesProcessor.
func NewFilesProcessor(files []string, processor *LineProcessor, writer io.Writer) *FilesProcessor {
	return &FilesProcessor{
		files:     files,
		processor: processor,
		writer:    writer,
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

		matched, err := fp.processor.Process(inputFile, fp.writer)
		if err != nil {
			return false, fmt.Errorf("error processing file %s: %w", file, err)
		}
		result = matched || result
	}
	return result, nil
}
