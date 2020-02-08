package dotfile

import (
	"bufio"
	"io"
)

func ReadSources(reader io.Reader, splitter bufio.SplitFunc) (sources []string, err error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitter)

	for scanner.Scan() {
		sources = append(sources, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return
}

const null = 0x0

var SplitOnNewlines bufio.SplitFunc = bufio.ScanLines

func SplitOnNullByte(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil // request more data
	}

	for i := 0; i < len(data); i++ {
		if data[i] == null {
			return i + 1, data[:i], nil
		}
	}
	if !atEOF {
		return 0, nil, nil
	}
	return 0, data, bufio.ErrFinalToken
}
