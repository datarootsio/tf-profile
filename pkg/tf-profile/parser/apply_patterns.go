package tfprofile

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

var (
	// All regexes that recognize interesting logs during the apply phase
	resourceName = `[a-zA-Z0-9_.["\]\/:]*` // Simplified regex but it will do

	resourceCreated         = fmt.Sprintf("%v: Creation complete after", resourceName)
	resourceCreationStarted = fmt.Sprintf("%v: Creating...", resourceName)
	resourceOperationFailed = fmt.Sprintf("with %v,", resourceName)

	resourceDestructionStarted = fmt.Sprintf("%v: Destroying...", resourceName)
	resourceDestroyed          = fmt.Sprintf("%v: Destruction complete after", resourceName)

	resourceModificationStarted = fmt.Sprintf("%v: Modifying...", resourceName)
	resourceModified            = fmt.Sprintf("%v: Modifications complete after", resourceName)
)

// Handle line that indicates creation of a resource was completed. E.g:
// resource: Creation complete after 1s [id=2023-04-09T18:17:33Z]
func parseResourceCreated(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceCreated, Line)
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
	createDuration := parseCreateDurationString(tokens2[0])

	// We know the resource and the duration, insert everything into the log
	log.SetTotalTime(resource, createDuration)
	log.SetAfterStatus(resource, Created)
	log.SetModificationCompletedEvent(resource, log.CurrentEvent)
	log.SetModificationCompletedIndex(resource, log.CurrentModificationEndedIndex)

	log.CurrentModificationEndedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates the creation of a resource was started. E.g:
// aws_ssm_parameter.bad2[2]: Creating...
func parseResourceCreationStarted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceCreationStarted, Line)
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
	log.SetOperation(tokens[0], Create)
	log.SetModificationStartedIndex(tokens[0], log.CurrentModificationStartedIndex)
	log.SetModificationStartedEvent(tokens[0], log.CurrentEvent)
	log.CurrentModificationStartedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates resource modifications failed. E.g:
// Error: creating SSM Parameter (/slash/at/end1/): ValidationException: Something something
//	 status code: 400, request id: 77765932-a8b2-48bf-abe2-71a151da56ea
//	 with aws_ssm_parameter.bad2[1],
// In practice we just detect the "with <resource_name>", as we only receive one line of context
func parseResourceCreationFailed(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceOperationFailed, Line)
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
	// TODO: dependin on the operation, Failed is not always correct. E.g. destroy fails => Created
	log.SetAfterStatus(resource, Failed)
	return true, nil
}

// Handle line that indicates the destruction of a resource was started. E.g:
// aws_ssm_parameter.bad2[2]: Destroying...
func parseResourceDestructionStarted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceDestructionStarted, Line)
	if !match {
		return false, nil
	}
	tokens := strings.Split(Line, ": Destroying...")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource deletion line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}

	// Knowing the resource whose deletion stared, insert everything in the log
	log.RegisterNewResource(tokens[0])
	log.SetOperation(tokens[0], Destroy)
	log.SetModificationCompletedEvent(tokens[0], log.CurrentEvent)
	log.SetModificationCompletedIndex(tokens[0], log.CurrentModificationEndedIndex)
	log.CurrentModificationStartedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates deletion of a resource was completed. E.g:
// resource: Destruction complete after 1s [id=2023-04-09T18:17:33Z]
func parseResourceDestroyed(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceDestroyed, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, ": Destruction complete after ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource destruction line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := tokens[0]

	// The next token will contain the create time (" Destruction complete after ...s [id=...]")
	tokens2 := strings.Split(tokens[1], " ")
	createDuration := parseCreateDurationString(tokens2[0])

	// We know the resource and the duration, insert everything into the log
	log.SetTotalTime(resource, createDuration)
	log.SetAfterStatus(resource, NotCreated)
	log.SetModificationCompletedEvent(resource, log.CurrentEvent)
	log.SetModificationCompletedIndex(resource, log.CurrentModificationEndedIndex)

	log.CurrentModificationEndedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates the destruction of a resource was started. E.g:
// aws_ssm_parameter.bad2[2]: Destroying...
func parseResourceModificationStarted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceModificationStarted, Line)
	if !match {
		return false, nil
	}
	tokens := strings.Split(Line, ": Modifying...")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource modification line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}

	// Knowing the resource whose modification stared, insert everything in the log
	log.RegisterNewResource(tokens[0])
	log.SetOperation(tokens[0], Modify)
	log.SetModificationStartedEvent(tokens[0], log.CurrentEvent)
	log.SetModificationStartedIndex(tokens[0], log.CurrentModificationStartedIndex)
	log.CurrentModificationStartedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}

// Handle line that indicates modification of a resource was completed. E.g:
// resource: Destruction complete after 1s [id=2023-04-09T18:17:33Z]
func parseResourceModified(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceModified, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, ": Modifications complete after ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource modification line: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := tokens[0]

	// The next token will contain the create time (" Modifications complete after ...s [id=...]")
	tokens2 := strings.Split(tokens[1], " ")
	if len(tokens2) < 2 {
		msg := fmt.Sprintf("Unable to parse duration: %v\n", tokens[1])
		return false, &LineParseError{Msg: msg}
	}
	Duration := parseCreateDurationString(tokens2[0])

	// We know the resource and the duration, insert everything into the log
	log.SetTotalTime(resource, Duration)
	log.SetAfterStatus(resource, Created)
	log.SetModificationCompletedEvent(resource, log.CurrentEvent)
	log.SetModificationCompletedIndex(resource, log.CurrentModificationEndedIndex)

	log.CurrentModificationEndedIndex += 1
	log.CurrentEvent += 1
	return true, nil
}
