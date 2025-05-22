package patt

import (
	"bufio"
	"io"
)

func PrintMatchingLines(filter LinesMatcher, reader io.Reader, writer io.Writer) (bool, error) {
	scanner := bufio.NewScanner(reader)
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()

	match := false
	for scanner.Scan() {
		line := scanner.Bytes()
		if filter.Match(line) {
			match = true
			_, err := bufferedWriter.Write(line)
			if err != nil {
				return false, err
			}
			err = bufferedWriter.WriteByte('\n')
			if err != nil {
				return false, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}
	return match, nil
}
