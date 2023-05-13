package tfprofile

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

var (
	ResourceName            = `[a-zA-Z0-9_.["\]\/:]*` // Simplified regex but it will do
	ResourceCreated         = fmt.Sprintf("%v: Creation complete after", ResourceName)
	ResourceCreationStarted = fmt.Sprintf("%v: Creating...", ResourceName)
	ResourceCreationFailed  = fmt.Sprintf("with %v,", ResourceName)
)

// Handle line that indicates creation of a resource was completed. E.g:
// resource: Creation complete after 1s [id=2023-04-09T18:17:33Z]
func ParseResourceCreated(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ResourceCreated, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, ": Creation complete after ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := tokens[0]

	// The next token will contain the create time (" Creation complete after ...s [id=...]")
	tokens2 := strings.Split(tokens[1], " ")
	if len(tokens2) < 2 {
		msg := fmt.Sprintf("Unable to parse creation duration: %v\n", tokens[1])
		return false, &LineParseError{Msg: msg}
	}
	createDuration := ParseCreateDurationString(tokens2[0])

	// We know the resource and the duration, insert everything into the log
	log.SetTotalTime(resource, createDuration)
	log.SetCreationStatus(resource, Created)
	log.SetCreationCompletedEvent(resource, log.CurrentEvent)
	log.SetCreationCompletedIndex(resource, log.CurrentModificationEndedIndex)

	log.CurrentModificationEndedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates the creation of a resource was started. E.g:
// aws_ssm_parameter.bad2[2]: Creating...
func ParseResourceCreationStarted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ResourceCreationStarted, Line)
	if !match {
		return false, nil
	}
	tokens := strings.Split(Line, ": Creating...")
	if len(tokens) < 2 || tokens[1] != "" {
		msg := fmt.Sprintf("Unable to parse resource creation line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}

	// Knowing the resource whose creation stared, insert everything in the log
	log.RegisterNewResource(tokens[0])
	log.CurrentModificationStartedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates resource modifications failed. E.g:
// Error: creating SSM Parameter (/slash/at/end1/): ValidationException: Something something
//	 status code: 400, request id: 77765932-a8b2-48bf-abe2-71a151da56ea
//	 with aws_ssm_parameter.bad2[1],
// In practice we just detect the "with <resource_name>", as we only receive one line of context
func ParseResourceCreationFailed(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ResourceCreationFailed, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "with ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse failure line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	Line = strings.TrimSpace(Line)
	resource := strings.Split(Line, "with ")[1] // Everything after 'with'
	resource = resource[:len(resource)-1]       // Remove comma at end

	// Knowing the resource whose modifications failed, insert everything in the log
	log.SetCreationStatus(resource, Failed)
	return true, nil
}
