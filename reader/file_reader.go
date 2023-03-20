package reader

import (
	"bufio"
	"os"
)

type FileReader struct {
	File string
}

func (r FileReader) Read() (*bufio.Scanner, error) {
	file, err := os.Open(r.File)
	if err != nil {
		return nil, err
	}

	return bufio.NewScanner(file), nil
}
