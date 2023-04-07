package tfprofile

import (
	"sort"
	"strings"
)

// Take a parsed log and aggregate resources created
// by the same `foreach` or `count` loop.
func Aggregate(log ParsedLog) (ParsedLog, error) {
	New := ParsedLog{resources: make(map[string]ResourceMetric)}

	// Collect all resource names in slice and sort
	ResourceNames := []string{}
	for k, _ := range log.resources {
		ResourceNames = append(ResourceNames, k)
	}
	sort.Strings(ResourceNames)

	// Loop over resources, collect resources to aggregate
	// and aggregate when coming across a new one.
	ToAgg := []string{}

	for _, name := range ResourceNames {
		// If the list is empty, start over.
		if len(ToAgg) == 0 {
			ToAgg = append(ToAgg, name)
		} else if canAggregate(name, last(ToAgg)) {
			// If this name is compatible with the list, add it
			ToAgg = append(ToAgg, name)
		} else {
			// Found incompatible resource, aggregate, add to log and start over
			AggName, AggMetric := aggregateResources(log, ToAgg)

			// Add aggregated resource to log we're building
			NewLog := New.resources
			NewLog[AggName] = AggMetric
			New.resources = NewLog

			// Start over with what the resource we just saw
			ToAgg = []string{name}
		}
	}

	// Aggregate leftovers in the list to get the final result
	if len(ToAgg) > 0 {
		AggName, AggMetric := aggregateResources(log, ToAgg)
		NewLog := New.resources
		NewLog[AggName] = AggMetric
		New.resources = NewLog
	}

	return New, nil
}

// Helper: return the last item in a list
func last(l []string) string {
	return l[len(l)-1]
}

// Given two resource names, returns true if they were created using
// `count` or `for_each`. For example: `resource[1]` and `resource[2]`
func canAggregate(Resource1 string, Resource2 string) bool {
	// If any resource doesn't end with "]", can't aggregate
	if Resource1[len(Resource1)-1] != ']' || Resource2[len(Resource2)-1] != ']' {
		return false
	}
	// Anything before the last "[" must be equal
	split1 := strings.Split(Resource1, "[")
	split2 := strings.Split(Resource2, "[")
	if len(split1) == 1 || len(split2) == 1 {
		return false
	}

	prefix1 := strings.Join(split1[:len(split1)-1], "[")
	prefix2 := strings.Join(split2[:len(split2)-1], "[")
	if prefix1 != prefix2 {
		return false
	}
	return true
}

// Given a log and resources names to aggregate, find an aggregated name and
// an aggregated ResourceMetric
func aggregateResources(log ParsedLog, resources []string) (string, ResourceMetric) {
	// Singleton, just return it
	if len(resources) == 1 {
		return resources[0], log.resources[resources[0]]
	}

	Metrics := []ResourceMetric{}
	for _, r := range resources {
		Metrics = append(Metrics, log.resources[r])
	}

	return aggregateResourceNames(resources...), aggregateResourceMetrics(Metrics...)
}

// Returns a new name for aggregated resource. For example:
// (module.x.resource[1], module.x.resource[2]) -> module.x.resource[*]
func aggregateResourceNames(names ...string) string {
	// Actually we only need to look at one item for now.
	split := strings.Split(names[0], "[")
	return strings.Join(split[:len(split)-1], "[") + "[*]"
}

// Aggregates a number of ResourceMetrics into one.
// After aggregating 'NumCalls' contains the number of input records.
// TotalTime contains the sum of individual apply times.
// CreationStartedIndex contains the *lowest* CreationStartedIndex of any record.
// CreationCompletedIndex contains the *highest* CreationStartedIndex of any record.
// CreationStatus can be any of "AllCreated", "AllStarted", "NoneStarted", "SomeFailed",
// "AllFailed"
func aggregateResourceMetrics(metrics ...ResourceMetric) ResourceMetric {
	NumCalls := len(metrics)
	TotalTime := float64(0)
	CreationStartedIndex := -1
	CreationCompletedIndex := -1

	AllCreatedB := true
	AllStartedB := true
	NoneStartedB := true
	SomeFailedB := false
	AllFailedB := true

	for _, metric := range metrics {
		TotalTime += metric.TotalTime
		if CreationStartedIndex == -1 {
			CreationStartedIndex = metric.CreationStartedIndex
		}
		CreationCompletedIndex = maxInt(CreationCompletedIndex, metric.CreationCompletedIndex)
		if metric.CreationStatus == Created {
			NoneStartedB = false
			AllFailedB = false
		}
		if metric.CreationStatus == Failed {
			AllCreatedB = false
			NoneStartedB = false
			SomeFailedB = true
		}
		if metric.CreationStatus == Started {
			AllCreatedB = false
			NoneStartedB = false
			AllFailedB = false
		}
		if metric.CreationStatus == NotStarted {
			AllCreatedB = false
			AllStartedB = false
			AllFailedB = false
		}
	}

	// FinalStatus should be the most interesting status
	// we can give based on the metrics seen
	var FinalStatus Status
	if AllCreatedB {
		FinalStatus = AllCreated
	} else if AllFailedB {
		FinalStatus = AllFailed
	} else if SomeFailedB {
		FinalStatus = SomeFailed
	} else if NoneStartedB {
		FinalStatus = NoneStarted
	} else if AllStartedB {
		FinalStatus = AllStarted
	} else {
		FinalStatus = SomeStarted
	}

	return ResourceMetric{
		NumCalls:               NumCalls,
		TotalTime:              TotalTime,
		CreationStartedIndex:   CreationStartedIndex,
		CreationCompletedIndex: CreationCompletedIndex,
		CreationStatus:         FinalStatus,
	}
}

func maxInt(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}
