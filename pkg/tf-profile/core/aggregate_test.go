package tfprofile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregateResourceMetricBasic(t *testing.T) {
	M1 := ResourceMetric{1, 2000, 0, 0, Created}
	M2 := ResourceMetric{1, 5000, 1, 1, Created}
	M3 := ResourceMetric{1, 1000, 2, 2, Created}

	Result := aggregateResourceMetrics(M1, M2, M3)
	Expected := ResourceMetric{3, 8000, 0, 2, AllCreated}
	assert.Equalf(t, Expected, Result, "Expected different result after aggregating.")
}

func AggStatus(In ...Status) Status {
	ResourceMetrics := []ResourceMetric{}
	for _, rm := range In {
		ResourceMetrics = append(ResourceMetrics, ResourceMetric{CreationStatus: rm})
	}
	return aggregateResourceMetrics(ResourceMetrics...).CreationStatus
}
func TestAggregateResourceMetricStatuses(t *testing.T) {
	Result := AggStatus(Failed, Failed, Failed)
	assert.Equal(t, AllFailed, Result)

	Result = AggStatus(Created, Failed, NotStarted)
	assert.Equal(t, SomeFailed, Result)

	Result = AggStatus(Started, Started, NotStarted)
	assert.Equal(t, SomeStarted, Result)

	Result = AggStatus(Created, Created, NotStarted)
	assert.Equal(t, SomeStarted, Result)

	Result = AggStatus(Created, Created, Failed)
	assert.Equal(t, SomeFailed, Result)
}

func TestCanAgg(t *testing.T) {
	assert.False(t, canAggregate("resource1", "resource2"))
	assert.False(t, canAggregate("module.x.r1", "module.y.r1"))
	assert.False(t, canAggregate("module.x[1].r1", "module.x[1].r2"))
	assert.False(t, canAggregate("module.x.r[1]", "module.y.r[2]"))
	assert.False(t, canAggregate("module.x.r1[\"abc\"]", "module.y.r1[\"def\"]"))

	assert.True(t, canAggregate("r1[1]", "r1[2]"))
	assert.True(t, canAggregate("r1[\"abc\"]", "r1[\"def\"]"))
	assert.True(t, canAggregate("module.x.r1[\"abc\"]", "module.x.r1[\"def\"]"))
	assert.True(t, canAggregate("module.x[\"a\"].r1[\"abc\"]", "module.x[\"a\"].r1[\"def\"]"))
	assert.True(t, canAggregate("r[1]", "r[\"a\"]")) // Edge case as they come from different loops...
}

// No aggregation possible
func TestNoAgg(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"resource1": ResourceMetric{},
			"resource2": ResourceMetric{},
			"resource3": ResourceMetric{},
		},
	}
	Result, err := Aggregate(In)
	assert.Nil(t, err)
	assert.Equal(t, In, Result) // Assert Result identical to input
}

func TestBasicAgg(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r1[1]": ResourceMetric{1, 1, 0, 0, Created},
			"r1[2]": ResourceMetric{1, 1, 1, 1, Created},
			"r1[3]": ResourceMetric{1, 1, 2, 2, Created},
		},
	}
	Out := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r1[*]": ResourceMetric{3, 3, 0, 2, AllCreated},
		},
	}
	Result, err := Aggregate(In)
	assert.Nil(t, err)
	assert.Equal(t, Out, Result)
}

func TestMixedAgg(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r1[1]":     ResourceMetric{1, 1, 0, 0, Created},
			"r1[2]":     ResourceMetric{1, 1, 1, 1, Created},
			"r1[3]":     ResourceMetric{1, 1, 2, 2, Created},
			"r2[\"a\"]": ResourceMetric{1, 1, 3, 3, Created},
			"r2[\"b\"]": ResourceMetric{1, 1, 4, 4, Created},
			"r3":        ResourceMetric{1, 1, 5, 5, Created},
			"r4":        ResourceMetric{1, 1, 6, 6, Created},
		},
	}
	Out := ParsedLog{
		Resources: map[string]ResourceMetric{
			"r1[*]": ResourceMetric{3, 3, 0, 2, AllCreated},
			"r2[*]": ResourceMetric{2, 2, 3, 4, AllCreated},
			"r3":    ResourceMetric{1, 1, 5, 5, Created},
			"r4":    ResourceMetric{1, 1, 6, 6, Created},
		},
	}
	Result, err := Aggregate(In)
	assert.Nil(t, err)
	assert.Equal(t, Out, Result)
}

func TestFullAgg(t *testing.T) {
	In := ParsedLog{
		Resources: map[string]ResourceMetric{
			// Can be aggregated on name
			"module.x.r[1]": ResourceMetric{1, 1, 0, 0, Created},
			"module.x.r[2]": ResourceMetric{1, 2, 0, 0, Created},
			"module.x.r[3]": ResourceMetric{1, 3, 0, 0, Created},
			// With nested modules
			"module.y[1].module.y[1].r[1]": ResourceMetric{1, 1, 0, 0, Created},
			"module.y[1].module.y[1].r[2]": ResourceMetric{1, 2, 0, 0, Created},
			"module.y[1].module.y[1].r[3]": ResourceMetric{1, 3, 0, 0, Created},
			// With CreationStartedIndex and CreationCompletedIndex
			"module.z[1].module.z[1].r[1]": ResourceMetric{1, 3, 2, 1, Created},
			"module.z[1].module.z[1].r[2]": ResourceMetric{1, 2, 3, 5, Created},
			"module.z[1].module.z[1].r[3]": ResourceMetric{1, 1, 5, 9, Created},
			// Mixed states
			"module.a[1].module.b[1].r[1]": ResourceMetric{1, 3, 2, 1, Created},
			"module.a[1].module.b[1].r[2]": ResourceMetric{1, 2, 3, 5, Failed},
			"module.a[1].module.b[1].r[3]": ResourceMetric{1, 1, 5, 9, Started},
			"module.a[1].module.b[1].r[4]": ResourceMetric{1, 1, 5, 9, NotStarted},
			// Not agg'able
			"random_resource":  ResourceMetric{1, 1, 5, 9, Created},
			"random_resource2": ResourceMetric{1, 1, 5, 9, Failed},
			"random_resource3": ResourceMetric{1, 1, 5, 9, NotStarted},
		},
	}
	Out := ParsedLog{
		Resources: map[string]ResourceMetric{
			"module.x.r[*]":                ResourceMetric{3, 6, 0, 0, AllCreated},
			"module.y[1].module.y[1].r[*]": ResourceMetric{3, 6, 0, 0, AllCreated},
			"module.z[1].module.z[1].r[*]": ResourceMetric{3, 6, 2, 9, AllCreated},
			"module.a[1].module.b[1].r[*]": ResourceMetric{4, 7, 2, 9, SomeFailed},
			"random_resource":              ResourceMetric{1, 1, 5, 9, Created},
			"random_resource2":             ResourceMetric{1, 1, 5, 9, Failed},
			"random_resource3":             ResourceMetric{1, 1, 5, 9, NotStarted},
		},
	}
	Result, err := Aggregate(In)
	assert.Nil(t, err)
	assert.Equal(t, Out, Result)
}
