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
	ResourceName    = `[a-zA-Z0-9_.["\]]*` // Simplified regex but it will do
	ResourceCreated = fmt.Sprintf("%v: Creation complete after", ResourceName)
)

type (
	LineParseError struct{ Msg string }

	// Data structure that holds all metrics for one particular resource
	ResourceMetric struct {
		NumCalls      int
		TotalTime     float64
		CreationIndex int // Resource was the Nth to start creation, not implemented
		CreatedIndex  int // Resource was the Nth to finish creation
	}

	// Parsing a log results in a map of resource names and their metrics
	ParsedLog = map[string]ResourceMetric
)

func (e *LineParseError) Error() string {
	return e.Msg
}

func Parse(file *bufio.Scanner, tee bool) (ParsedLog, error) {
	num_created := 0

	tflog := make(ParsedLog)
	for file.Scan() {
		line := file.Text()
		if tee {
			fmt.Println(line)
		}

		match, _ := regexp.MatchString(ResourceCreated, line)
		if match {
			resource, time, err := ParseResourceCreated(line)
			if err != nil {
				msg := `This line was detected to contain resource creation, 
				            but tf-profile is unable to parse it!`
				return nil, &LineParseError{msg}
			}
			InsertResourceMetric(tflog, resource, time, num_created)
			num_created += 1
		}
	}

	return tflog, nil
}

// Insert a new ResourceMetric into the log
func InsertResourceMetric(log ParsedLog, resource string, duration float64, idx int) {
	metric, exists := (log)[resource]

	// We have seen this resource before
	if exists {
		(metric).NumCalls += 1
		(metric).TotalTime += 1
	} else {
		// Add new resource with some default values
		(log)[resource] = ResourceMetric{
			NumCalls:      1,
			TotalTime:     duration,
			CreationIndex: -1, // Not implemented
			CreatedIndex:  idx,
		}
	}
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
