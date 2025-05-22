package patt

import (
	"bufio"
	"io"
)

func PrintMatchingLines(filter LinesMatcher, reader io.Reader, writer io.Writer) (bool, error) {
	scanner := bufio.NewScanner(reader)

	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	bufferedWriter := bufio.NewWriterSize(writer, 256*1024)
	defer bufferedWriter.Flush()

	match := false

	for scanner.Scan() {
		line := scanner.Bytes()
		if filter.Match(line) {
			match = true

			lineWithNewline := append(line, '\n')
			if _, err := bufferedWriter.Write(lineWithNewline); err != nil {
				return false, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return match, nil
}
