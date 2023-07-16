package tfprofile

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/readers"
)

type Stat struct {
	name  string
	value string
}

func Filter(args []string) error {
	var file *bufio.Scanner
	var regex string
	var err error

	file, regex, err = parseArgs(args)
	if err != nil {
		return err
	}

	logs := FilterLogs(file, regex)
	printList(logs)
	return nil
}

// Read a file line by line and return only lines that contain information
// regarding the resources we're filtering to. This includes: plan logs,
// errors logs and (otherwise) any lines that match the resource regex.
func FilterLogs(file *bufio.Scanner, regex string) []string {
	output := []string{}

	// If we're currently seeing plan logs, collect the lines into a
	// buffer until detecting the end, then add all at once.
	var CollectingPlan = false
	var PlanBuffer = []string{}

	// If we're seeing errors logs, collect lines into a buffer until
	// the end of the message, then print add at once.
	var CollectingError = false
	var ErrorBuffer = []string{}

	for file.Scan() {
		line := RemoveTerminalFormatting(file.Text())

		if CollectingPlan && isEndOfPlan(line) {
			// End of plan, print the buffer and reset
			PlanBuffer = append(PlanBuffer, line)
			output = append(output, PlanBuffer...)
			output = append(output, "") // Empty line to make it look nice.

			PlanBuffer = []string{}
			CollectingPlan = false
		} else if CollectingPlan {
			// Continuation of plan
			PlanBuffer = append(PlanBuffer, line)
		} else if isStartOfPlan(line, regex) {
			// Detected start of plan, start new buffer
			CollectingPlan = true
			PlanBuffer = append(PlanBuffer, line)
		} else if CollectingError && isEndOfError(line) {
			// End of error, print buffer and reset
			ErrorBuffer = append(ErrorBuffer, line)
			output = append(output, ErrorBuffer...)

			ErrorBuffer = []string{}
			CollectingError = false
		} else if CollectingError && errorDoesntMatchResource(line, regex) {
			// Discard buffer, error is not about an interesting resource
			ErrorBuffer = []string{}
			CollectingError = false
		} else if CollectingError {
			// Continuation of error
			ErrorBuffer = append(ErrorBuffer, line)
		} else if isStartOfError(line) {
			// Start new error
			CollectingError = true
			ErrorBuffer = append(ErrorBuffer, line)
		} else {
			// No plan or error but a resource is mentioned: print the line
			match, _ := regexp.MatchString(regex, line)
			if match {
				output = append(output, line)
			}
		}
	}
	return output
}

// Make using regex to specify Terraform resources easier by allowing things
// that are not valid regex but make sense  intuitively.
// For example: `module.*.my_resource[*]` can not directly be evaluated because
// `.` and `[]` are regex constructs, `*` is used as `.*`. To allow this type of
// "natural" querying, we replace `*` with `.*` and escape the following characters:
// ".[]".
func cleanRegex(regex string) string {
	out := strings.ReplaceAll(regex, `.`, `\.`)
	out = strings.ReplaceAll(out, `*`, `.*`)
	out = strings.ReplaceAll(out, `]`, `\]`)
	out = strings.ReplaceAll(out, `[`, `\[`)
	return out
}

func parseArgs(args []string) (*bufio.Scanner, string, error) {
	var err error
	var file *bufio.Scanner
	var regex string

	if len(args) == 2 {
		file, err = FileReader{File: args[1]}.Read()
		regex = args[0]

		if file == nil {
			return nil, "", fmt.Errorf("File could not be read (%v)\n", args[1])
		}
	} else if len(args) == 1 {
		file, err = StdinReader{}.Read()
		regex = args[0]
	} else if err != nil {
		return nil, "", err
	} else {
		return nil, "", fmt.Errorf(
			"Filter command requires one or two arguments, %v were given!\n", len(args))
	}

	return file, cleanRegex(regex), nil
}

// The start of a plan block can be identified by various sentences, such as:
// "X will be created", "X will replaced", etc...
func isStartOfPlan(line string, resource string) bool {
	patterns := []string{
		fmt.Sprintf("%v is tainted, so must be replaced", resource),
		fmt.Sprintf("%v will be created", resource),
		fmt.Sprintf("%v will be replaced, as requested", resource),
		fmt.Sprintf("%v will be destroyed", resource),
		fmt.Sprintf("%v will be updated in-place", resource),
		fmt.Sprintf("%v must be replaced", resource),
	}

	// Check if any of the patterns match
	for _, p := range patterns {
		match, _ := regexp.MatchString(p, line)
		if match {
			return true
		}
	}
	return false
}

func isEndOfPlan(line string) bool {
	return line == "    }"
}

func printList(b []string) {
	for _, line := range b {
		fmt.Println(line)
	}
}

func isStartOfError(line string) bool {
	return strings.HasPrefix(line, "Error: ")
}

func isEndOfError(line string) bool {
	pattern := `[0-9]+: (resource|data) ".*" ".*" {`
	match, _ := regexp.MatchString(pattern, line)
	return match
}

// Return true when we are detecting an error but it does not relate
// to a resource (specified by regex).
func errorDoesntMatchResource(line string, resource string) bool {
	// Interesting lines start with "  with " and and with ","
	if strings.HasPrefix(line, "  with ") && strings.HasSuffix(line, ",") {
		pattern := fmt.Sprintf("  with %v,", resource)
		match, _ := regexp.MatchString(pattern, line)
		return !match
	}
	return false
}
