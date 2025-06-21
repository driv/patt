package patt

import (
	"bufio"
	"io"
)

type LineProcessor struct {
	reader          *bufio.Scanner
	writer          *bufio.Writer
	keepNonMatching bool
}

func NewLineProcessor(r io.Reader, w io.Writer, keepNonMatching bool) *LineProcessor {
	return &LineProcessor{
		reader:          bufio.NewScanner(r),
		writer:          bufio.NewWriter(w),
		keepNonMatching: keepNonMatching,
	}
}

func (p *LineProcessor) ProcessLines(replacer LineReplacer) (bool, error) {
	defer p.writer.Flush()

	match := false
	for p.reader.Scan() {
		line := p.reader.Bytes()
		matched, err := p.processLine(line, replacer)
		if err != nil {
			return false, err
		}
		if matched {
			match = true
		}
	}

	if err := p.reader.Err(); err != nil {
		return false, err
	}
	return match, nil
}

func (p *LineProcessor) processLine(line []byte, replacer LineReplacer) (bool, error) {
	var out []byte
	var matched bool
	out, matched, b, err := newFunction(replacer, line, out, matched, p)
	if err != nil {
		return b, err
	}
	if p.keepNonMatching || matched {
		if err := p.writeLine(out); err != nil {
			return false, err
		}
	}
	return matched, nil
}

func newFunction(replacer LineReplacer, line []byte, out []byte, matched bool, p *LineProcessor) ([]byte, bool, bool, error) {
	if replacer.Match(line) {
		replaced, err := replacer.Replace(line)
		if err != nil {
			return nil, false, false, err
		}
		out = replaced
		matched = true
	} else if p.keepNonMatching {
		out = line
	}
	return out, matched, false, nil
}

func (p *LineProcessor) writeLine(line []byte) error {
	if _, err := p.writer.Write(line); err != nil {
		return err
	}
	return p.writer.WriteByte('\n')
}
