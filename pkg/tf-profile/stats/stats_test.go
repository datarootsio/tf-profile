package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestBasicStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": {NumCalls: 1, AfterStatus: Created},
			"b": {NumCalls: 1, AfterStatus: Created},
			"c": {NumCalls: 1, AfterStatus: Created},
			"d": {NumCalls: 1, AfterStatus: Created},
		},
	}
	Out := getBasicStats(In)
	assert.Equal(t, 1, len(Out))
	assert.Equal(t, "Number of resources in configuration", Out[0].name)
	assert.Equal(t, "4", Out[0].value)
}

func TestTimeStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": {NumCalls: 1, TotalTime: 1000, AfterStatus: Created},
			"b": {NumCalls: 1, TotalTime: 2000, AfterStatus: Created},
			"c": {NumCalls: 1, TotalTime: 3000, AfterStatus: Created},
			"d": {NumCalls: 1, TotalTime: 59000, AfterStatus: Created},
		},
	}
	Out := getTimeStats(In)

	assert.Equal(t, 3, len(Out))
	assert.Equal(t, "Cumulative duration", Out[0].name)
	assert.Equal(t, "Longest apply time", Out[1].name)
	assert.Equal(t, "Longest apply resource", Out[2].name)

	assert.Equal(t, "1m5s", Out[0].value)
	assert.Equal(t, "59s", Out[1].value)
	assert.Equal(t, "d", Out[2].value)
}

func TestStatusStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": {NumCalls: 1, AfterStatus: Created},
			"b": {NumCalls: 1, AfterStatus: Failed},
			"c": {NumCalls: 1, AfterStatus: Failed},
			"d": {NumCalls: 1, AfterStatus: NotCreated},
		},
	}
	Out := getAfterStatusStats(In)

	Expected := []Stat{
		{"Resources in state Created", "1"},
		{"Resources in state Failed", "2"},
		{"Resources in state NotCreated", "1"},
	}
	assert.Equal(t, Expected, Out)
}

func TestModuleStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r.test":                                            {NumCalls: 1, AfterStatus: Created},
			"module.test1.resource.test1":                       {NumCalls: 1, AfterStatus: Created},
			"module.test1.resource.test2":                       {NumCalls: 1, AfterStatus: Created},
			"module.test1.resource.test3":                       {NumCalls: 1, AfterStatus: Created},
			"module.test2.resource.test1":                       {NumCalls: 1, AfterStatus: Created},
			"module.test2.resource.test2":                       {NumCalls: 1, AfterStatus: Created},
			"module.a.module.b.module.c.module.d.resource.test": {NumCalls: 1, AfterStatus: Created},
		},
	}
	Out := getModuleStats(In)
	Expected := []Stat{
		{"Number of top-level modules", "3"},
		{"Largest top-level module", "module.test1"},
		{"Size of largest top-level module", "3"},
		{"Deepest module", "module.a.module.b.module.c.module.d"},
		{"Deepest module depth", "4"},
		{"Largest leaf module", "module.test1"},
		{"Size of largest leaf module", "3"},
	}
	assert.Equal(t, Expected, Out)
}

func TestFullStats(t *testing.T) {
	err := Stats([]string{"../../../test/aggregate.log"}, false, true)
	assert.Nil(t, err)
	err = Stats([]string{"../../../test/multiple_resources.log"}, false, true)
	assert.Nil(t, err)
	err = Stats([]string{"../../../test/null_resources.log"}, false, true)
	assert.Nil(t, err)

	err = Stats([]string{"does-not-exist"}, false, true)
	assert.NotNil(t, err)
}
