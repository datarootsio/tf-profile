package tfprofile

import "fmt"

const (
	// Status for individual resources
	// NotStarted Status = 0
	// Started    Status = 1
	NoneStatus Status = -1 // Internal only
	Unknown    Status = 0
	NotCreated Status = 1
	Created    Status = 2
	Failed     Status = 3
	Tainted    Status = 4
	// For aggregated resources
	Multiple Status = 4

	// Operation types
	NoneOp     Operation = -1 // Internal only
	None       Operation = 0  // Default, "not seen", can occur when upstream resources fail
	Create     Operation = 1
	Modify     Operation = 2
	Replace    Operation = 3
	Destroy    Operation = 4
	MultipleOp Operation = 5
)

type (
	Status    int
	Operation int

	// Data structure that holds all metrics for one particular resource
	ResourceMetric struct {
		NumCalls  int
		TotalTime float64
		// Resource was the Nth to start creation.
		ModificationStartedIndex int
		// Resource was the Nth to finish creation
		ModificationCompletedIndex int
		// (Global) event index of when creation started. As this is a global event,
		// it can be compared chronologically with a ModificationCompletedEvent.
		ModificationStartedEvent int
		// (Global) event index of when creation finished. As this is a global event,
		// it can be compared chronologically with a ModificationStartedEvent.
		ModificationCompletedEvent int // (Global) event index of when creation finished
		// Inferred status before the TF run
		BeforeStatus Status
		// Status after the TF run
		AfterStatus Status
		// Expected status as planned by TF
		DesiredStatus Status
		// Operation to perform to go from BeforeStatus to DesiredStatus
		Operation Operation
	}

	// Parsing a log results in a map of resource names and their metrics
	ParsedLog struct {
		// Indices to keep track of progress during parse
		CurrentModificationStartedIndex int
		CurrentModificationEndedIndex   int
		CurrentEvent                    int
		// Stage information
		ContainsRefresh bool
		ContainsPlan    bool
		ContainsApply   bool
		// Resources detected
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

func (log ParsedLog) SetModificationStartedIndex(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.ModificationStartedIndex = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetModificationCompletedIndex(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.ModificationCompletedIndex = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetModificationStartedEvent(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.ModificationStartedEvent = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetModificationCompletedEvent(Resource string, Idx int) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.ModificationCompletedEvent = Idx
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetAfterStatus(Resource string, Status Status) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.AfterStatus = Status
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetBeforeStatus(Resource string, Status Status) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.BeforeStatus = Status
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetDesiredStatus(Resource string, Status Status) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	metric.DesiredStatus = Status
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) SetOperation(Resource string, Op Operation) error {
	metric, found := log.Resources[Resource]
	if found == false {
		return &ResourceNotFoundError{Resource}
	}
	// If Operation was Destroy before, overwrite Create with Replace
	if metric.Operation == Destroy && Op == Create {
		metric.Operation = Replace
	} else {
		metric.Operation = Op
	}
	log.Resources[Resource] = metric
	return nil
}

func (log ParsedLog) RegisterNewResource(Resource string) {
	_, found := (log.Resources)[Resource]
	if found {
		return
	}
	(log.Resources)[Resource] = ResourceMetric{
		NumCalls:                   1,
		TotalTime:                  -1, // Not finished yet, will be overwritten
		ModificationStartedIndex:   log.CurrentModificationStartedIndex,
		ModificationCompletedIndex: -1, // Not finished yet, will be overwritten
		ModificationStartedEvent:   log.CurrentEvent,
		ModificationCompletedEvent: -1, // Not finished yet, will be overwritten

		AfterStatus: NotCreated, // Not finished yet, will be overwritten
	}
}

func (s Status) String() string {
	switch s {
	case NotCreated:
		return "NotCreated"
	case Created:
		return "Created"
	case Failed:
		return "Failed"
	case Unknown:
		return "Unknown"
	case Tainted:
		return "Tainted"
	default:
		return fmt.Sprintf("%d (unknown)", int(s))
	}
}

func (s Operation) String() string {
	switch s {
	case Destroy:
		return "Destroy"
	case Create:
		return "Create"
	case Modify:
		return "Modify"
	case Replace:
		return "Replace"
	case MultipleOp:
		return "Multiple"
	case None:
		return "None recorded"
	default:
		return fmt.Sprintf("%d (unknown)", int(s))
	}
}
