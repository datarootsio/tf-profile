package readers

import "bufio"

// A Reader creates a Scanner to read a log line by line
type Reader interface {
	OpenFile() *bufio.Scanner
}
