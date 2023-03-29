package tfprofile

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicParseResourceCreated(t *testing.T) {
	line := "x[0]: Creation complete after 10s [id=xyz]"
	name, time, err := ParseResourceCreated(line)

	if err != nil || name != "x[0]" || time != 10000 {
		t.Fatalf("%v - %v - %v\n", name, time, err)
	}

	line2 := "x[\"a\"]: Creation complete after 5s [id=xyz]"
	name, time, err = ParseResourceCreated(line2)

	if err != nil || name != "x[\"a\"]" || time != 5000 {
		t.Fatalf("%v - %v - %v\n", name, time, err)
	}
}

func TestParseResourceCreatedMinutes(t *testing.T) {
	_, time, _ := ParseResourceCreated("x: Creation complete after 5m30s [id=xyz]")
	expected_time := float64((5*60 + 30) * 1000)

	if time != expected_time {
		t.Fatalf("Expected %v, got %v\n", expected_time, time)
	}
}

func TestFullParse(t *testing.T) {
	file, _ := os.Open("../../test/multiple_resources.log")
	s := bufio.NewScanner(file)

	log, err := Parse(s, false)
	assert.Nil(t, err)

	metrics, ok := log["time_sleep.count[9]"]
	assert.True(t, ok)

	expected := ResourceMetric{
		NumCalls:      1,
		TotalTime:     10000,
		CreationIndex: -1, // Not implemented
		CreatedIndex:  12,
	}
	if metrics != expected {
		t.Fatalf("Expected %v, got %v\n", expected, metrics)
	}

	metrics2 := log["time_sleep.for_each[\"a\"]"]
	expected2 := ResourceMetric{
		NumCalls:      1,
		TotalTime:     1000,
		CreationIndex: -1, // Not implemented
		CreatedIndex:  1,
	}
	if metrics2 != expected2 {
		t.Fatalf("Expected %v, got %v\n", expected2, metrics2)
	}
}
