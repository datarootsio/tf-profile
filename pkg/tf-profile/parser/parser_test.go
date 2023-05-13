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
		NumCalls:               1,
		TotalTime:              10000,
		CreationStartedIndex:   10,
		CreationCompletedIndex: 12,
		CreationStartedEvent:   11,
		CreationCompletedEvent: 26,
		CreationStatus:         Created,
	}
	if metrics != expected {
		t.Fatalf("Expected %v, got %v\n", expected, metrics)
	}

	metrics2 := log.Resources["time_sleep.for_each_a"]
	expected2 := ResourceMetric{
		NumCalls:               1,
		TotalTime:              1000,
		CreationStartedIndex:   5,
		CreationCompletedIndex: 1,
		CreationStartedEvent:   5,
		CreationCompletedEvent: 12,
		CreationStatus:         Created,
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
	assert.Equal(t, metrics.CreationStatus, Created)

	metrics, exists = log.Resources["aws_ssm_parameter.good"]
	assert.True(t, exists)
	assert.Equal(t, metrics.CreationStatus, Created)

	metrics, exists = log.Resources["aws_ssm_parameter.bad2[1]"]
	assert.True(t, exists)
	assert.Equal(t, metrics.CreationStatus, Failed)

	metrics, exists = log.Resources["aws_ssm_parameter.bad"]
	assert.True(t, exists)
	assert.Equal(t, metrics.CreationStatus, Failed)
}

// func TestSpecialResourceName(t *testing.T) {
// 	// Name contains special characters: :/_
// 	line := `module.eks.module.eks_managed_node_group["initial"].aws_iam_role_policy_attachment.this["arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"]: Creation complete after 0s [id=initial-eks-node-group-20230430082916737400000001-20230430082917966400000003]"`
// 	name, duration, err := ParseResourceCreated(line)

// 	assert.Nil(t, err)
// 	assert.Equal(t, `module.eks.module.eks_managed_node_group["initial"].aws_iam_role_policy_attachment.this["arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"]`, name)
// 	assert.Equal(t, float64(0), duration)
// }

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
