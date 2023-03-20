package printer

import (
	"bufio"
	"os"
	"testing"

	"github.com/QuintenBruynseraede/tf-profile/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseSortSpec(t *testing.T) {
	p1 := parseSortSpec("key=asc")
	assert.Equal(t, len(p1), 1, "Expected one item after parsing.")
	assert.Equal(t, p1[0].col, "key")
	assert.Equal(t, p1[0].order, "asc")

	p2 := parseSortSpec("a=asc,b=desc,c=asc")
	expected := []SortSpecItem{
		SortSpecItem{"a", "asc"},
		SortSpecItem{"b", "desc"},
		SortSpecItem{"c", "asc"},
	}
	assert.Equalf(t, p2, expected, "Expected %v after parsing, got %v\n", p2, expected)
}

func TestSort(t *testing.T) {
	file, _ := os.Open("../test_files/multiple_resources.log")
	s := bufio.NewScanner(file)
	log, err := parser.Parse(s, false)
	assert.Nil(t, err)

	sorted := Sort(&log, "tot_time=asc,idx_created=asc")
	expected := []string{
		"time_sleep.count[0]",
		"time_sleep.for_each[\"a\"]",
		"time_sleep.count[1]",
		"time_sleep.for_each[\"c\"]",
		"time_sleep.count[2]",
		"time_sleep.for_each[\"d\"]",
		"time_sleep.for_each[\"b\"]",
		"time_sleep.count[3]",
		"time_sleep.count[4]",
		"time_sleep.count[5]",
		"time_sleep.count[6]",
		"time_sleep.count[7]",
		"time_sleep.count[8]",
		"time_sleep.count[9]",
	}
	assert.Equal(t, sorted, expected)

	sorted2 := Sort(&log, "tot_time=desc,idx_created=desc")
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], sorted2[len(expected)-i-1])
	}
	Table(&log, "tot_time=asc,idx_created=asc")
}
