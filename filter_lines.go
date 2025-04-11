package patt

import (
	"bufio"
	"io"
)

func MatchLines(filter LinesMatcher, reader io.Reader, writer io.Writer) (bool, error) {
	bufferedReader := bufio.NewReader(reader)
	match := false
	cont := true
	for cont {
		line, err := bufferedReader.ReadBytes(newLine)
		if err == io.EOF {
			cont = false
			if len(line) == 0 {
				break
			}
		} else if err != nil {
			return false, err
		} else {
			line = line[:len(line)-1]
		}

		if !filter.Match(line) {
			continue
		}
		match = true
		if err := WriteLine(writer, line); err != nil {
			return false, err
		}
	}
	return match, nil
}

const newLine = '\n'

func WriteLine(writer io.Writer, line []byte) error {
	_, err := writer.Write(line)
	_, _ = writer.Write([]byte{newLine})
	return err
}
