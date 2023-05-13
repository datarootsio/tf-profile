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

var RefreshParsers = []ParseFunction{}
var PlanParsers = []ParseFunction{
	ParseStartPlan,
}
var ApplyParsers = []ParseFunction{
	ParseResourceCreationStarted,
	ParseResourceCreated,
	ParseResourceCreationFailed,
	ParseResourceDestructionStarted,
	ParseResourceDestroyed,
	ParseResourceModificationStarted,
	ParseResourceModified,
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

		// Apply refresh parsers until one modifies the log
		for _, f := range RefreshParsers {
			modified, err := f(line, &tflog)
			if err != nil {
				return ParsedLog{}, err
			}
			if modified {
				tflog.ContainsRefresh = true
				break
			}
		}

		// Apply plan parsers until one modifies the log
		for _, f := range PlanParsers {
			modified, err := f(line, &tflog)
			if err != nil {
				return ParsedLog{}, err
			}
			if modified {
				tflog.ContainsPlan = true
				break
			}
		}

		// Apply apply parsers until one modifies the log.
		for _, f := range ApplyParsers {
			modified, err := f(line, &tflog)
			if err != nil {
				return ParsedLog{}, err
			}
			if modified {
				tflog.ContainsApply = true
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
