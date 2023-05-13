package tfprofile

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

// Parse a refresh line and records the resource in the log.
func refreshParser(Line string, log *ParsedLog) (bool, error) {
	regex := fmt.Sprintf("%v: Refreshing state...", resourceName)
	match, _ := regexp.MatchString(regex, Line)
	if !match {
		return false, nil
	}
	tokens := strings.Split(Line, ": Refreshing state...")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}

	// Knowing the resource whose creation stared, insert everything in the log
	resource := tokens[0]
	log.RegisterNewResource(resource)
	log.SetModificationStartedEvent(resource, -1)
	log.SetModificationStartedIndex(resource, -1)

	return true, nil
}
