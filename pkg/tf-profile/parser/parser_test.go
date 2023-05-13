package tfprofile

import (
	"bufio"
	"os"
	"strings"
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"

	"github.com/stretchr/testify/assert"
)

func TestFullParse(t *testing.T) {
	file, _ := os.Open("../../../test/multiple_resources.log")
	s := bufio.NewScanner(file)

	log, err := Parse(s, false)
	assert.Nil(t, err)

	metrics, ok := log.Resources["time_sleep.count_9"]
	assert.True(t, ok)

	expected := ResourceMetric{
		NumCalls:                   1,
		TotalTime:                  10000,
		ModificationStartedIndex:   10,
		ModificationCompletedIndex: 12,
		ModificationStartedEvent:   11,
		ModificationCompletedEvent: 26,
		BeforeStatus:               Created,
		DesiredStatus:              Created,
		AfterStatus:                Created,
		Operation:                  Create,
	}
	if metrics != expected {
		t.Fatalf("Expected %v, got %v\n", expected, metrics)
	}

	metrics2 := log.Resources["time_sleep.for_each_a"]
	expected2 := ResourceMetric{
		NumCalls:                   1,
		TotalTime:                  1000,
		ModificationStartedIndex:   5,
		ModificationCompletedIndex: 1,
		ModificationStartedEvent:   5,
		ModificationCompletedEvent: 12,
		BeforeStatus:               Created,
		DesiredStatus:              Created,
		AfterStatus:                Created,
		Operation:                  Create,
	}
	if metrics2 != expected2 {
		t.Fatalf("Expected %v, got %v\n", expected2, metrics2)
	}
}

func TestFailureParse(t *testing.T) {
	file, _ := os.Open("../../../test/failures.log")
	s := bufio.NewScanner(file)

	log, err := Parse(s, false)
	assert.Nil(t, err)

	metrics, exists := log.Resources["aws_ssm_parameter.good2[0]"]
	assert.True(t, exists)
	assert.Equal(t, metrics.AfterStatus, Created)

	metrics, exists = log.Resources["aws_ssm_parameter.good"]
	assert.True(t, exists)
	assert.Equal(t, metrics.AfterStatus, Created)

	metrics, exists = log.Resources["aws_ssm_parameter.bad2[1]"]
	assert.True(t, exists)
	assert.Equal(t, metrics.AfterStatus, Failed)

	metrics, exists = log.Resources["aws_ssm_parameter.bad"]
	assert.True(t, exists)
	assert.Equal(t, metrics.AfterStatus, Failed)
}

func TestParserSanityCheck(t *testing.T) {
	Files, err := os.ReadDir("../../../test")
	assert.Nil(t, err)

	// Sanity check: all *.log files must be parseable
	for _, File := range Files {
		if strings.Contains(File.Name(), ".log") {
			file, _ := os.Open("../../../test/" + File.Name())
			s := bufio.NewScanner(file)

			_, err := Parse(s, false)
			assert.Nil(t, err)
		}
	}
}
