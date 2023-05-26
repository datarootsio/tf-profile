package tfprofile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/aggregate"
	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/parser"

	"github.com/stretchr/testify/assert"
)

func TestFailureParse(t *testing.T) {
	Files, err := os.ReadDir("../../../test")
	assert.Nil(t, err)

	// Sanity check: all *.log files must be graph-able
	for _, File := range Files {
		if strings.Contains(File.Name(), ".log") {
			err := Graph([]string{"../../../test/" + File.Name()}, 1000, 600, "tf-profile-graph.png", true)
			assert.Nil(t, err)
		}
	}

	err = Graph([]string{"../../../test/does-not-exist"}, 1000, 600, "tf-profile-graph.png", true)
	assert.NotNil(t, err)
	err = Graph([]string{"../../../test/failures.log"}, -1, -1, "tf-profile-graph.png", true)
	assert.NotNil(t, err)
}

func TestPlotOutput(t *testing.T) {
	file, _ := os.Open("../../../test/failures.log")
	s := bufio.NewScanner(file)

	log, _ := Parse(s, false)
	log, _ = Aggregate(log)

	out, err := printGNUPlotOutput(log, 1000, 600, "tf-profile-graph.png")

	assert.Nil(t, err)
	fmt.Println(out)
	assert.Contains(t, out, `aws\\\_ssm\\\_parameter.good2[*] 7 11 Created`)
	assert.Contains(t, out, `aws\\\_ssm\\\_parameter.bad 5 -1 Failed`)
	assert.Contains(t, out, `aws\\\_ssm\\\_parameter.bad2[*] 3 -1 Failed`)
	assert.Contains(t, out, `aws\\\_ssm\\\_parameter.good 0 8 Created`)

}
