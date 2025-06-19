package patt

import (
	"bufio"
	"io"
)

type LineProcessor struct {
	reader *bufio.Scanner
	writer *bufio.Writer
}

func NewLineProcessor(r io.Reader, w io.Writer) *LineProcessor {
	return &LineProcessor{
		reader: bufio.NewScanner(r),
		writer: bufio.NewWriter(w),
	}
}

func PrintMatchingLines(filter LinesMatcher, reader io.Reader, writer io.Writer) (bool, error) {
	replacer := matchFilter{LinesMatcher: filter}
	processor := NewLineProcessor(reader, writer)
	return processor.processLines(replacer)
}

func PrintLines(replacer LineReplacer, reader io.Reader, writer io.Writer) (bool, error) {
	processor := NewLineProcessor(reader, writer)
	return processor.processLines(replacer)
}

type matchFilter struct {
	LinesMatcher
}

func (mf matchFilter) Replace(line []byte) ([]byte, error) {
	return line, nil
}

func (p *LineProcessor) processLines(replacer LineReplacer) (bool, error) {
	defer p.writer.Flush()

	match := false
	for p.reader.Scan() {
		line := p.reader.Bytes()
		if !replacer.Match(line) {
			continue
		}
		match = true
		line, err := replacer.Replace(line)
		if err != nil {
			return false, err
		}
		if err := p.writeLine(line); err != nil {
			return false, err
		}
	}

	if err := p.reader.Err(); err != nil {
		return false, err
	}
	return match, nil
}

func (p *LineProcessor) writeLine(line []byte) error {
	if _, err := p.writer.Write(line); err != nil {
		return err
	}
	return p.writer.WriteByte('\n')
}
