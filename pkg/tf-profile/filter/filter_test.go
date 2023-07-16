package tfprofile

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	err := Filter([]string{"aws_ssm_parameter.*", "../../../test/failures.log"})
	assert.Nil(t, err)

	err = Filter([]string{})
	assert.NotNil(t, err)

	err = Filter([]string{"*", "non/existing/file.txt"})
	assert.NotNil(t, err)

	err = Filter([]string{"1", "2", "3"})
	assert.NotNil(t, err)

	err = Filter([]string{"*"})
	assert.Nil(t, err)
}

func TestFilterSanityCheck(t *testing.T) {
	Files, err := os.ReadDir("../../../test")
	assert.Nil(t, err)

	// Sanity check: some output when filtering to resource .*
	for _, File := range Files {
		if strings.Contains(File.Name(), ".log") {
			file, _ := os.Open("../../../test/" + File.Name())
			s := bufio.NewScanner(file)

			out := FilterLogs(s, ".*")
			assert.True(t, len(out) > 0)
		}
	}
}

func TestFullyQualifiedResourceFilter(t *testing.T) {
	file, _ := os.Open("../../../test/null_resources.log")
	s := bufio.NewScanner(file)
	regex := cleanRegex("null_resource.next")
	out := FilterLogs(s, regex)

	assert.Contains(t, out, `  # null_resource.next will be created`)
	assert.Contains(t, out, `  + resource "null_resource" "next" {`)
	assert.Contains(t, out, `      + id = (known after apply)`)
	assert.Contains(t, out, `    }`)
	assert.Contains(t, out, `null_resource.next: Creating...`)
	assert.Contains(t, out, `null_resource.next: Creation complete after 0s [id=1766626192520212902]`)
}

func TestBasicWildcardFilter(t *testing.T) {
	file, _ := os.Open("../../../test/null_resources.log")
	s := bufio.NewScanner(file)
	regex := cleanRegex("null_resource*")
	out := FilterLogs(s, regex)

	assert.Contains(t, out, `  # null_resource.next will be created`)
	assert.Contains(t, out, `  # null_resource.previous will be created`)
	assert.Contains(t, out, `  + resource "null_resource" "next" {`)
	assert.Contains(t, out, `  + resource "null_resource" "previous" {`)
	assert.Contains(t, out, `null_resource.previous: Creation complete after 0s [id=5144705655797302376]`)
	assert.Contains(t, out, `null_resource.next: Creation complete after 0s [id=1766626192520212902]`)
}

func TestFilterWithError(t *testing.T) {
	file, _ := os.Open("../../../test/failures.log")
	s := bufio.NewScanner(file)
	regex := cleanRegex("aws_ssm_parameter.bad2*")
	out := FilterLogs(s, regex)

	assert.Contains(t, out, `  # aws_ssm_parameter.bad2[0] will be created`)
	assert.Contains(t, out, `  # aws_ssm_parameter.bad2[1] will be created`)
	assert.Contains(t, out, `  # aws_ssm_parameter.bad2[2] will be created`)

	assert.Contains(t, out, `Error: creating SSM Parameter (/slash/at/end2/): ValidationException: Parameter name must not end with slash.`)
	assert.Contains(t, out, `Error: creating SSM Parameter (/slash/at/end1/): ValidationException: Parameter name must not end with slash.`)
	assert.Contains(t, out, `Error: creating SSM Parameter (/slash/at/end0/): ValidationException: Parameter name must not end with slash.`)

	assert.Contains(t, out, `  with aws_ssm_parameter.bad2[0],`)
	assert.Contains(t, out, `  with aws_ssm_parameter.bad2[1],`)
	assert.Contains(t, out, `  with aws_ssm_parameter.bad2[2],`)

	// Doesn't contain anything related to aws_ssm_parameter.bad
	assert.NotContains(t, out, `  with aws_ssm_parameter.bad,`)
	assert.NotContains(t, out, `Error: creating SSM Parameter (/slash/at/end/): ValidationException: Parameter name must not end with slash.`)
}

func TestCleanRegex(t *testing.T) {
	var m = map[string]string{
		`*`:             `.*`,
		`x.y`:           `x\.y`,
		`module.*.x[*]`: `module\..*\.x\[.*\]`,
	}

	for in, out := range m {
		assert.Equal(t, cleanRegex(in), out)
	}
}

func TestPatterns(t *testing.T) {
	assert.True(t, isStartOfPlan("x is tainted, so must be replaced", "x"))
	assert.True(t, isStartOfPlan("x will be created", "x"))
	assert.True(t, isStartOfPlan("x will be replaced, as requested", "x"))
	assert.True(t, isStartOfPlan("x will be destroyed", "x"))
	assert.True(t, isStartOfPlan("x will be updated in-place", "x"))
	assert.True(t, isStartOfPlan("x must be replaced", "x"))

	assert.True(t, isStartOfError("Error: something the provider returns"))
	assert.True(t, isEndOfError(`12: resource "x" "y" {`))
	assert.True(t, isEndOfError(`1234: data "x" "y" {`))
	assert.False(t, isEndOfError(`12: something "x" "x" {`))
	assert.False(t, isEndOfError(`abcd: resource "x" "x" {`))
	assert.False(t, isEndOfError(`12: resource "x" "x"`))

	assert.True(t, errorDoesntMatchResource(`  with y,`, `x`))
	assert.True(t, errorDoesntMatchResource(`  with aws_ssm_parameter,`, `azure*`))
}
