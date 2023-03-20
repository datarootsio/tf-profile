package reader

import (
	"bufio"
	"os"
)

type StdinReader struct{}

func (r StdinReader) Read() (*bufio.Scanner, error) {
	return bufio.NewScanner(os.Stdin), nil
}
