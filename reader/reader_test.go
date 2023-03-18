package reader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReader(t *testing.T) {
	inputfile := "../test_files/test_file.txt"
	file := FileReader{File: inputfile}.Read()

	file.Scan()
	assert.Equal(t, file.Text(), "Used to test reader module")
	file.Scan()
	assert.Equal(t, file.Text(), "Another line")
	file.Scan()
	assert.Equal(t, file.Text(), "Another line2")
	assert.Equal(t, false, file.Scan())
	assert.Equal(t, file.Text(), "")
}
