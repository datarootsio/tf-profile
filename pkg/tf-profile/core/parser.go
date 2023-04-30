package tfprofile

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	ResourceName            = `[a-zA-Z0-9_.["\]\/:]*` // Simplified regex but it will do
	ResourceCreated         = fmt.Sprintf("%v: Creation complete after", ResourceName)
	ResourceCreationStarted = fmt.Sprintf("%v: Creating...", ResourceName)
	ResourceCreationFailed  = "Error: "
)

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
	Status         int
	LineParseError struct{ Msg string }

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
		creationStartedCount  int
		createdCompletedCount int
		Resources             map[string]ResourceMetric
	}
)

func (e *LineParseError) Error() string {
	return e.Msg
}

func Parse(file *bufio.Scanner, tee bool) (ParsedLog, error) {
	CreationStarted := 0
	CreationCompleted := 0
	EventIndex := 0 // Any start or ending of a creation/modification/deletion is an event
	// In case a resource update fails, the resource name comes a couple of lines after
	// the error. This flag is true when we are looking for the resource after an error.
	FailureSeen := false

	tflog := ParsedLog{0, 0, make(map[string]ResourceMetric)}

	for file.Scan() {
		line := file.Text()
		if tee {
			fmt.Println(line)
		}

		match, _ := regexp.MatchString(ResourceCreated, line)
		if match {
			resource, time, err := ParseResourceCreated(line)
			if err != nil {
				msg := fmt.Sprintf(`Resource name detected (CreationEnded), but unable to parse line "%v"`, line)
				return ParsedLog{}, &LineParseError{msg}
			}
			FinishResourceCreation(tflog, resource, time, CreationCompleted, EventIndex)
			CreationCompleted += 1
			EventIndex += 1
		}

		match, _ = regexp.MatchString(ResourceCreationStarted, line)
		if match {
			resource, err := ParseResourceCreationStarted(line)
			if err != nil {
				msg := fmt.Sprintf(`Resource name detected (CreationStarted), but unable to parse line "%v"`, line)
				return ParsedLog{}, &LineParseError{msg}
			}
			InsertResourceMetric(tflog, resource, CreationStarted, EventIndex)
			CreationStarted += 1
			EventIndex += 1
		}

		match, _ = regexp.MatchString(ResourceCreationFailed, line)
		if match {
			FailureSeen = true // Start counting
		}
		if FailureSeen == true {
			if strings.Contains(line, "with ") && strings.Contains(line, ",") {
				resource, err := ParseCreationFailed(line)
				if err != nil {
					msg := fmt.Sprintf(`Line contains failed resource name, we can't parse it: "%v"`, line)
					return ParsedLog{}, &LineParseError{msg}
				}
				MarkAsFailed(tflog, resource)
				FailureSeen = false
			}
		}
	}

	return tflog, nil
}

// Insert a new ResourceMetric into the log
func InsertResourceMetric(log ParsedLog, resource string, CreationIndex int, EventIndex int) {
	(log.Resources)[resource] = ResourceMetric{
		NumCalls:               1,
		TotalTime:              -1, // Not finished yet, will be overwritten
		CreationStartedIndex:   CreationIndex,
		CreationCompletedIndex: -1, // Not finished yet, will be overwritten
		CreationStartedEvent:   EventIndex,
		CreationCompletedEvent: -1, // Not finished yet, will be overwritten
		CreationStatus:         Started,
	}
}

// When creation is done, update the record in the log
func FinishResourceCreation(Log ParsedLog, Resource string, Duration float64, Idx int, Event int) error {
	record, exists := Log.Resources[Resource]

	if exists == false {
		msg := fmt.Sprintf("Resource %v finished creation, but start was not recorded!\n", Resource)
		return &LineParseError{msg}
	}

	record.CreationCompletedIndex = Idx
	record.TotalTime = Duration
	record.CreationStatus = Created
	record.CreationCompletedEvent = Event
	Log.Resources[Resource] = record

	return nil
}

// Mark a resource in the log as Failed
func MarkAsFailed(log ParsedLog, resource string) error {
	record, exists := log.Resources[resource]
	if exists == false {
		msg := fmt.Sprintf("Creation of resource %v failed, but its was not seen before!\n", resource)
		return &LineParseError{msg}
	}
	record.CreationStatus = Failed
	log.Resources[resource] = record
	return nil
}

// Parse one line containing resource creation text
func ParseResourceCreated(line string) (string, float64, error) {
	tokens := strings.Split(line, ": Creation complete after ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", line)
		return "", -1, &LineParseError{msg}
	}
	resource := tokens[0]

	// The next token will contain the create time (" Creation complete after ...s [id=...]")
	tokens2 := strings.Split(tokens[1], " ")
	if len(tokens2) < 2 {
		msg := fmt.Sprintf("Unable to parse creation duration: %v\n", tokens[1])
		return "", -1, &LineParseError{msg}
	}
	create_duration := ParseCreateDurationString(tokens2[0])
	return resource, create_duration, nil
}

// Parse one line containing resource creation text
func ParseResourceCreationStarted(line string) (string, error) {
	tokens := strings.Split(line, ": Creating...")
	if len(tokens) < 2 || tokens[1] != "" {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", line)
		return "", &LineParseError{msg}
	}
	return tokens[0], nil
}

// Parse one line containing resource failure ("with <resource-name>,")
func ParseCreationFailed(line string) (string, error) {
	tokens := strings.Split(line, "with ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse failure line: %v\n", line)
		return "", &LineParseError{msg}
	}
	line = strings.TrimSpace(line)
	resource := strings.Split(line, "with ")[1] // Everything after 'with'
	return resource[:len(resource)-1], nil
}

// Convert a create duration string into milliseconds
func ParseCreateDurationString(in string) float64 {
	// Q: what's the formatting when > 1hr?
	// For now handle two cases: "1m10s" and "10s"
	if strings.Contains(in, "m") {
		split := strings.Split(in, "m")
		mins, err1 := strconv.Atoi(split[0])
		seconds, err2 := strconv.Atoi(strings.TrimSuffix(split[1], "s"))

		if err1 != nil || err2 != nil {
			log.Fatal("Unable to parse resource create duration.")
		}

		return float64(1000.0 * (60*mins + seconds))
	} else {
		seconds, err := strconv.Atoi(strings.TrimSuffix(in, "s"))
		if err != nil {
			log.Fatal("Unable to parse resource create duration.")
		}
		return float64(1000.0 * seconds)
	}
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
