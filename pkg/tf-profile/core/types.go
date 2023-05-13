package tfprofile

import "fmt"

const (
	// For individual resources
	NotStarted Status = 0
	Started    Status = 1
	Created    Status = 2
	Failed     Status = 3
	// For aggregated resources
	SomeStarted Status = 4
	AllStarted  Status = 5
	NoneStarted Status = 6
	SomeFailed  Status = 7
	AllFailed   Status = 8
	AllCreated  Status = 9
)

type (
	Status int

	// Data structure that holds all metrics for one particular resource
	ResourceMetric struct {
		NumCalls  int
		TotalTime float64
		// Resource was the Nth to start creation.
		CreationStartedIndex int
		// Resource was the Nth to finish creation
		CreationCompletedIndex int
		// (Global) event index of when creation started. As this is a global event,
		// it can be compared chronologically with a CreationCompletedEvent.
		CreationStartedEvent int
		// (Global) event index of when creation finished. As this is a global event,
		// it can be compared chronologically with a CreationStartedEvent.
		CreationCompletedEvent int // (Global) event index of when creation finished
		CreationStatus         Status
	}

	// Parsing a log results in a map of resource names and their metrics
	ParsedLog struct {
		Resources map[string]ResourceMetric
	}
)

func (log ParsedLog) SetNumCalls(Resource string, NumCalls int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.NumCalls = NumCalls
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetTotalTime(Resource string, TotalTime float64) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.TotalTime = TotalTime
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetCreationStartedIndex(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.CreationStartedIndex = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetCreationCompletedIndex(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.CreationCompletedIndex = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetCreationStartedEvent(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.CreationStartedEvent = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetCreationCompletedEvent(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.CreationCompletedEvent = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetCreationStatus(Resource string, Status Status) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.CreationStatus = Status
	log.Resources[Resource] = metric
	return nil
}

func (s Status) String() string {
	switch s {
	case NotStarted:
		return "NotStarted"
	case Started:
		return "Started"
	case Created:
		return "Created"
	case Failed:
		return "Failed"
	case SomeStarted:
		return "SomeStarted"
	case AllStarted:
		return "AllStarted"
	case NoneStarted:
		return "NoneStarted"
	case SomeFailed:
		return "SomeFailed"
	case AllFailed:
		return "AllFailed"
	case AllCreated:
		return "AllCreated"
	default:
		return fmt.Sprintf("%d (unknown)", int(s))
	}
}
