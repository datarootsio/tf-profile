package tfprofile

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

type parseFunction = func(Line string, log *ParsedLog) (bool, error)

var RefreshParsers = []parseFunction{
	refreshParser,
}
var PlanParsers = []parseFunction{
	parseStartPlan,
	parsePlanTainted,
	parsePlanExplicitReplace,
	parsePlanWillBeDestroyed,
	parsePlanWillBeModified,
	parsePlanForcedReplace,
	parsePlanWillBeCreated,
}
var ApplyParsers = []parseFunction{
	parseResourceCreationStarted,
	parseResourceCreated,
	parseResourceCreationFailed,
	parseResourceDestructionStarted,
	parseResourceDestroyed,
	parseResourceModificationStarted,
	parseResourceModified,
}

// Parse a Terraform log into a ParsedLog object. This function will
// pass line by line over the file, apply parse functions (see above)
// until one of them recognizes the line and extracts information. In
// that case the line is considered "handled" and the next one is scanned
// Possible optimization here: since Terraform has distinct refresh,
// plan, apply phases we could skip parse functions of previous phases.
func Parse(file *bufio.Scanner, tee bool) (ParsedLog, error) {
	// regex to detect ANSI terminal formatting directives (https://stackoverflow.com/a/14693789)
	re := regexp.MustCompile(`(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]`)

	tflog := ParsedLog{Resources: map[string]ResourceMetric{}}

	for file.Scan() {
		line := file.Text()
		line = re.ReplaceAllString(line, "")

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
func parseCreateDurationString(in string) float64 {
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
