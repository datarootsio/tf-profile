package tfprofile

import (
	"bufio"
	"os"
	"testing"

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
	file, _ := os.Open("../../../test/multiple_resources.log")
	s := bufio.NewScanner(file)
	log, err := Parse(s, false)
	assert.Nil(t, err)

	sorted := Sort(log, "tot_time=asc,idx_created=asc")
	expected := []string{
		"time_sleep.count_0",
		"time_sleep.for_each_a",
		"time_sleep.count_1",
		"time_sleep.for_each_c",
		"time_sleep.count_2",
		"time_sleep.for_each_d",
		"time_sleep.for_each_b",
		"time_sleep.count_3",
		"time_sleep.count_4",
		"time_sleep.count_5",
		"time_sleep.count_6",
		"time_sleep.count_7",
		"time_sleep.count_8",
		"time_sleep.count_9",
	}
	assert.Equal(t, sorted, expected)

	sorted2 := Sort(log, "tot_time=desc,idx_created=desc")
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], sorted2[len(expected)-i-1])
	}
}
