package tfprofile

import "bufio"

// A Reader creates a Scanner to read a log line by line
type Reader interface {
	Read() (*bufio.Scanner, error)
}
