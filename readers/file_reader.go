package readers

import (
	"bufio"
	"log"
	"os"
)

type FileReader struct {
	File string
}

func (r FileReader) Read() *bufio.Scanner {
	file, err := os.Open(r.File)
	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewScanner(file)
}
