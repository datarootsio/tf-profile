package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestBasicStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": ResourceMetric{1, 0, 0, 0, Created},
			"b": ResourceMetric{1, 0, 0, 0, Created},
			"c": ResourceMetric{1, 0, 0, 0, Created},
			"d": ResourceMetric{1, 0, 0, 0, Created},
		},
	}
	Out := GetBasicStats(In)
	assert.Equal(t, 1, len(Out))
	assert.Equal(t, "Number of resources created", Out[0].name)
	assert.Equal(t, "4", Out[0].value)
}

func TestTimeStats(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"a": ResourceMetric{1, 1000, 0, 0, Created},
			"b": ResourceMetric{1, 2000, 0, 0, Created},
			"c": ResourceMetric{1, 3000, 0, 0, Created},
			"d": ResourceMetric{1, 59000, 0, 0, Created},
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
			"a": ResourceMetric{1, 0, 0, 0, Created},
			"b": ResourceMetric{1, 0, 0, 0, Failed},
			"c": ResourceMetric{1, 0, 0, 0, Failed},
			"d": ResourceMetric{1, 0, 0, 0, AllCreated},
			"e": ResourceMetric{1, 0, 0, 0, SomeFailed},
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
			"r.test":                                            ResourceMetric{1, 0, 0, 0, Created},
			"module.test1.resource.test1":                       ResourceMetric{1, 0, 0, 0, Created},
			"module.test1.resource.test2":                       ResourceMetric{1, 0, 0, 0, Created},
			"module.test2.resource.test1":                       ResourceMetric{1, 0, 0, 0, Created},
			"module.test2.resource.test2":                       ResourceMetric{1, 0, 0, 0, Created},
			"module.a.module.b.module.c.module.d.resource.test": ResourceMetric{1, 0, 0, 0, Created},
		},
	}
	Out := GetModuleStats(In)
	Expected := []Stat{
		Stat{"Number of top-level modules", "3"},
		Stat{"Largest top-level module", "module.test1"},
		Stat{"Size of largest top-level module", "2"},
		Stat{"Deepest module", "module.a.module.b.module.c.module.d"},
		Stat{"Deepest module depth", "4"},
		Stat{"Largest leaf module", "module.test1"},
		Stat{"Size of largest leaf module", "2"},
	}
	assert.Equal(t, Expected, Out)
}
