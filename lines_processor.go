package patt

import (
	"bufio"
	"io"
)

type LineProcessor struct {
	keepNonMatching bool
	replacer        LineReplacer
}

func NewLineProcessor(replacer LineReplacer, keepNonMatching bool) *LineProcessor {
	return &LineProcessor{
		keepNonMatching: keepNonMatching,
		replacer:        replacer,
	}
}

func (p *LineProcessor) Process(r io.Reader, w io.Writer) (bool, error) {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	var match bool
	for scanner.Scan() {
		line := scanner.Bytes()
		if p.replacer.Match(line) {
			line = p.replacer.Replace(line)
			match = true
		} else if !p.keepNonMatching {
			continue
		}
		if err := writeLine(writer, line); err != nil {
			return false, err
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}
	return match, nil
}

func writeLine(w *bufio.Writer, line []byte) error {
	if _, err := w.Write(line); err != nil {
		return err
	}
	return w.WriteByte('\n')
}
