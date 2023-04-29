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
	ResourceName            = `[a-zA-Z0-9_.["\]]*` // Simplified regex but it will do
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
		NumCalls               int
		TotalTime              float64
		CreationStartedIndex   int // Resource was the Nth to start creation, not implemented
		CreationCompletedIndex int // Resource was the Nth to finish creation
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
	// In case a resource update fails, the resource name comes three lines after
	// the error. This counter keeps track of when we can read the resource name
	LineSinceFailure := -1

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
				msg := `This line was detected to contain resource creation, but tf-profile is unable to parse it!`
				return ParsedLog{}, &LineParseError{msg}
			}
			FinishResourceCreation(tflog, resource, time, CreationCompleted)
			CreationCompleted += 1
		}

		match, _ = regexp.MatchString(ResourceCreationStarted, line)
		if match {
			resource, err := ParseResourceCreationStarted(line)
			if err != nil {
				msg := `This line was detected to contain resource creation, but tf-profile is unable to parse it!`
				return ParsedLog{}, &LineParseError{msg}
			}
			InsertResourceMetric(tflog, resource, CreationStarted)
			CreationStarted += 1
		}

		match, _ = regexp.MatchString(ResourceCreationFailed, line)
		if match {
			LineSinceFailure = 0 // Start counting
		}
		if LineSinceFailure != -1 {
			LineSinceFailure += 1
			if LineSinceFailure == 4 {
				resource, err := ParseCreationFailed(line)
				if err != nil {
					msg := `This line was detected to contain a resource name, but tf-profile is unable to parse it!`
					return ParsedLog{}, &LineParseError{msg}
				}
				MarkAsFailed(tflog, resource)
				LineSinceFailure = -1 // Stop line counting
			}
		}
	}

	return tflog, nil
}

// Insert a new ResourceMetric into the log
func InsertResourceMetric(log ParsedLog, resource string, idx int) {
	(log.Resources)[resource] = ResourceMetric{
		NumCalls:               1,
		TotalTime:              -1, // Not finished yet, will be overwritten
		CreationStartedIndex:   idx,
		CreationCompletedIndex: -1, // Not finished yet, will be overwritten
		CreationStatus:         Started,
	}
}

// When creation is done, update the record in the log
func FinishResourceCreation(log ParsedLog, resource string, duration float64, idx int) error {
	record, exists := log.Resources[resource]

	if exists == false {
		msg := fmt.Sprintf("Resource %v finished creation, but start was not recorded!\n", resource)
		return &LineParseError{msg}
	}

	record.CreationCompletedIndex = idx
	record.TotalTime = duration
	record.CreationStatus = Created
	log.Resources[resource] = record

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
	tokens := strings.Split(line, ":")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", line)
		return "", -1, &LineParseError{msg}
	}
	resource := tokens[0]

	// The next token will contain the create time (" Creation complete after ...s [id=...]")
	tokens2 := strings.Split(tokens[1], " [id=")
	if len(tokens2) < 2 {
		msg := fmt.Sprintf("Unable to parse creation duration: %v\n", tokens[1])
		return "", -1, &LineParseError{msg}
	}
	create_duration := ParseCreateDurationString(tokens2[0][25:])
	return resource, create_duration, nil
}

// Parse one line containing resource creation text
func ParseResourceCreationStarted(line string) (string, error) {
	tokens := strings.Split(line, ":")
	if len(tokens) < 2 {
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
