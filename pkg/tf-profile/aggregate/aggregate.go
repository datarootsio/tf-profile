package tfprofile

import (
	"sort"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

// Take a parsed log and aggregate resources created
// by the same `foreach` or `count` loop.
func Aggregate(log ParsedLog) (ParsedLog, error) {
	New := ParsedLog{Resources: make(map[string]ResourceMetric)}

	// Collect all resource names in slice and sort
	ResourceNames := []string{}
	for k, _ := range log.Resources {
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
			NewLog := New.Resources
			NewLog[AggName] = AggMetric
			New.Resources = NewLog

			// Start over with what the resource we just saw
			ToAgg = []string{name}
		}
	}

	// Aggregate leftovers in the list to get the final result
	if len(ToAgg) > 0 {
		AggName, AggMetric := aggregateResources(log, ToAgg)
		NewLog := New.Resources
		NewLog[AggName] = AggMetric
		New.Resources = NewLog
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
		return resources[0], log.Resources[resources[0]]
	}

	Metrics := []ResourceMetric{}
	for _, r := range resources {
		Metrics = append(Metrics, log.Resources[r])
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
// ModificationStartedIndex contains the *lowest* ModificationStartedIndex of any record.
// ModificationCompletedIndex contains the *highest* ModificationStartedIndex of any record.
// AfterStatus can be any of "Created", "Failed", "NotCreated", "Multiple" or "Unknown"
func aggregateResourceMetrics(metrics ...ResourceMetric) ResourceMetric {
	NumCalls := len(metrics)
	TotalTime := float64(0)
	ModificationStartedIndex := -1
	ModificationCompletedIndex := -1
	ModificationStartedEvent := -1
	ModificationCompletedEvent := -1

	BeforeStatus := NoneStatus
	AfterStatus := NoneStatus
	DesiredStatus := NoneStatus
	Operation := NoneOp

	for _, metric := range metrics {
		TotalTime += metric.TotalTime

		// For ModificationStartedIndex and ModificationStartedEvent, take the first one we see
		if ModificationStartedIndex == -1 {
			ModificationStartedIndex = metric.ModificationStartedIndex
		}
		if ModificationStartedEvent == -1 {
			ModificationStartedEvent = metric.ModificationStartedEvent
		}

		// For ModificationCompletedIndex and ModificationCompletedEvent, take the maximum
		ModificationCompletedIndex = maxInt(ModificationCompletedIndex, metric.ModificationCompletedIndex)
		ModificationCompletedEvent = maxInt(ModificationCompletedEvent, metric.ModificationCompletedEvent)

		// Calculate aggregated statuses:
		// - if all statuses are equal to X, the result will be X
		// - if multiple statuses are seen, the result will be "Multiple"
		if BeforeStatus == NoneStatus {
			BeforeStatus = metric.BeforeStatus
		}
		if AfterStatus == NoneStatus {
			AfterStatus = metric.AfterStatus
		}
		if DesiredStatus == NoneStatus {
			DesiredStatus = metric.DesiredStatus
		}
		if Operation == NoneOp {
			Operation = metric.Operation
		}

		if BeforeStatus != metric.BeforeStatus {
			BeforeStatus = Multiple
		}
		if AfterStatus != metric.AfterStatus {
			AfterStatus = Multiple
		}
		if DesiredStatus != metric.DesiredStatus {
			DesiredStatus = Multiple
		}
		if Operation != metric.Operation {
			Operation = MultipleOp
		}

	}

	return ResourceMetric{
		NumCalls:                   NumCalls,
		TotalTime:                  TotalTime,
		ModificationStartedIndex:   ModificationStartedIndex,
		ModificationCompletedIndex: ModificationCompletedIndex,
		ModificationStartedEvent:   ModificationStartedEvent,
		ModificationCompletedEvent: ModificationCompletedEvent,
		BeforeStatus:               BeforeStatus,
		AfterStatus:                AfterStatus,
		DesiredStatus:              DesiredStatus,
		Operation:                  Operation,
	}
}

func maxInt(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}
