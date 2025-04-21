package patt

import (
	"bufio"
	"io"
	"sync"
)

func MatchLines(filter LinesMatcher, reader io.Reader, writer io.Writer) (bool, error) {

	scanner := bufio.NewScanner(reader)
	match := false
	bufferedWriter := bufio.NewWriterSize(writer, 512*1000)
	defer bufferedWriter.Flush()

	lineChan := make(chan []byte, 100)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for line := range lineChan {
			bufferedWriter.Write(line)
			bufferedWriter.WriteByte('\n')
		}
	}()

	for scanner.Scan() {
		line := scanner.Bytes()
		if filter.Match(line) {
			match = true
			lineChan <- line
		}
	}

	close(lineChan)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return match, nil
}
