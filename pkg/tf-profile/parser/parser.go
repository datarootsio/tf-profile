package tfprofile

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

type ParseFunction = func(Line string, log *ParsedLog) (bool, error)

var ParseFunctions = []ParseFunction{
	ParseResourceCreated,
	ParseResourceCreationStarted,
	ParseResourceCreationFailed,
}

func Parse(file *bufio.Scanner, tee bool) (ParsedLog, error) {
	// CreationStarted := 0
	// EventIndex := 0 // Any start or ending of a creation/modification/deletion is an event
	// In case a resource update fails, the resource name comes a couple of lines after
	// the error. This flag is true when we are looking for the resource after an error.
	// FailureSeen := false

	// tflog := ParsedLog{make(map[string]ResourceMetric)}
	tflog := ParsedLog{Resources: map[string]ResourceMetric{}}

	for file.Scan() {
		line := file.Text()
		if tee {
			fmt.Println(line)
		}

		// Apply parse functions until one modifies the log.
		// In that case, we consider the line handled and go to the next one.
		for _, f := range ParseFunctions {
			modified, err := f(line, &tflog)
			if err != nil {
				return ParsedLog{}, err
			}
			if modified {
				break
			}
		}
	}

	return tflog, nil
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
