package patt

import (
	"bufio"
	"io"
)

type LineProcessor struct {
	reader          *bufio.Scanner
	writer          *bufio.Writer
	keepNonMatching bool
	replacer        LineReplacer
}

func NewLineProcessor(r io.Reader, w io.Writer, replacer LineReplacer, keepNonMatching bool) *LineProcessor {
	return &LineProcessor{
		reader:          bufio.NewScanner(r),
		writer:          bufio.NewWriter(w),
		keepNonMatching: keepNonMatching,
		replacer:        replacer,
	}
}

func (p *LineProcessor) Process() (bool, error) {
	defer p.writer.Flush()

	var match bool
	for p.reader.Scan() {
		line := p.reader.Bytes()
		if p.replacer.Match(line) {
			line = p.replacer.Replace(line)
			match = true
		} else if !p.keepNonMatching {
			continue
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
