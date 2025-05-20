package patt

import (
	"bufio"
	"io"
)

func ReplaceLines(filter LineReplacer, reader io.Reader, writer io.Writer) (bool, error) {
	scanner := bufio.NewScanner(reader)
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()

	match := false
	for scanner.Scan() {
		line := scanner.Bytes()
		if !filter.Match(line) {
			continue
		}
		line, err := filter.Replace(line)
		if err != nil {
			return false, err
		}
		match = true
		bufferedWriter.Write(line)
		bufferedWriter.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}
	return match, nil
}
