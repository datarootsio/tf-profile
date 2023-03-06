package readers

// A Reader reads a terraform log and calculates profiling metrics
// for all resources found in the logs
type Reader interface {
	ReadFile() []ResourceMetric
}

// Data structure that holds all metrics for one particular resource
type ResourceMetric struct {
	Resource  string
	NumCalls  int64
	TotalTime float64
	AvgTime   float64
}
