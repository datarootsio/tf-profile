package tfprofile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileReader(t *testing.T) {
	inputfile := "../../../test/test_file.txt"
	file, _ := FileReader{File: inputfile}.Read()

	file.Scan()
	assert.Equal(t, file.Text(), "Used to test reader module")
	file.Scan()
	assert.Equal(t, file.Text(), "Another line")
	file.Scan()
	assert.Equal(t, file.Text(), "Another line2")
	assert.Equal(t, false, file.Scan())
	assert.Equal(t, file.Text(), "")
}

func TestNonExistentFile(t *testing.T) {
	file, err := FileReader{File: "does-not-exist"}.Read()
	assert.Nil(t, file)
	assert.NotNil(t, err)
}
