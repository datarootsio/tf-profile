package readers

import (
	"bufio"
	"fmt"
	"os"
)

type StdinReader struct {
	Tee bool
}

func (r StdinReader) ReadFile() []ResourceMetric {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		if r.Tee {
			fmt.Println(line)
		}
		// Do something with line
	}
	fmt.Println("Read file.")
	return []ResourceMetric{}
}
