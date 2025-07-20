package patt

import (
	"bufio"
	"context"
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

const contextCheckInterval = 1000

func (p *LineProcessor) Process(ctx context.Context, r io.Reader, w io.Writer) (bool, error) {
	scanner := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	var match bool
	lines := 0
	for scanner.Scan() {
		lines++
		if lines%contextCheckInterval == 0 {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			default:
			}
		}
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
