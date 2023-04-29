package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestBasicStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"b": ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"c": ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"d": ResourceMetric{NumCalls: 1, CreationStatus: Created},
		},
	}
	Out := GetBasicStats(In)
	assert.Equal(t, 1, len(Out))
	assert.Equal(t, "Number of resources in configuration", Out[0].name)
	assert.Equal(t, "4", Out[0].value)
}

func TestTimeStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": ResourceMetric{NumCalls: 1, TotalTime: 1000, CreationStatus: Created},
			"b": ResourceMetric{NumCalls: 1, TotalTime: 2000, CreationStatus: Created},
			"c": ResourceMetric{NumCalls: 1, TotalTime: 3000, CreationStatus: Created},
			"d": ResourceMetric{NumCalls: 1, TotalTime: 59000, CreationStatus: Created},
		},
	}
	Out := GetTimeStats(In)

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
			"a": ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"b": ResourceMetric{NumCalls: 1, CreationStatus: Failed},
			"c": ResourceMetric{NumCalls: 1, CreationStatus: Failed},
			"d": ResourceMetric{NumCalls: 1, CreationStatus: AllCreated},
			"e": ResourceMetric{NumCalls: 1, CreationStatus: SomeFailed},
		},
	}
	Out := GetCreationStatusStats(In)

	Expected := []Stat{
		Stat{"No. resources in state AllCreated", "1"},
		Stat{"No. resources in state Created", "1"},
		Stat{"No. resources in state Failed", "2"},
		Stat{"No. resources in state SomeFailed", "1"},
	}
	assert.Equal(t, Expected, Out)
}

func TestModuleStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r.test":                                            ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.test1.resource.test1":                       ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.test1.resource.test2":                       ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.test1.resource.test3":                       ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.test2.resource.test1":                       ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.test2.resource.test2":                       ResourceMetric{NumCalls: 1, CreationStatus: Created},
			"module.a.module.b.module.c.module.d.resource.test": ResourceMetric{NumCalls: 1, CreationStatus: Created},
		},
	}
	Out := GetModuleStats(In)
	Expected := []Stat{
		Stat{"Number of top-level modules", "3"},
		Stat{"Largest top-level module", "module.test1"},
		Stat{"Size of largest top-level module", "3"},
		Stat{"Deepest module", "module.a.module.b.module.c.module.d"},
		Stat{"Deepest module depth", "4"},
		Stat{"Largest leaf module", "module.test1"},
		Stat{"Size of largest leaf module", "3"},
	}
	assert.Equal(t, Expected, Out)
}

func TestFullStats(t *testing.T) {
	err := Stats([]string{"../../../test/aggregate.log"}, false)
	assert.Nil(t, err)
	err = Stats([]string{"../../../test/multiple_resources.log"}, false)
	assert.Nil(t, err)
	err = Stats([]string{"../../../test/null_resources.log"}, false)
	assert.Nil(t, err)

	err = Stats([]string{"does-not-exist"}, false)
	assert.NotNil(t, err)
}
