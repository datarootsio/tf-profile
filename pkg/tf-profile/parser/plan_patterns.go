package tfprofile

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

var (
	startPlan       = "Terraform will perform the following actions:"
	isTainted       = fmt.Sprintf("%v is tainted, so must be replaced", resourceName)
	willBeCreated   = fmt.Sprintf("%v will be created", resourceName)
	explicitReplace = fmt.Sprintf("%v will be replaced, as requested", resourceName)
	willBeDestroyed = fmt.Sprintf("%v will be destroyed", resourceName)
	willBeModified  = fmt.Sprintf("%v will be updated in-place", resourceName)
	forcedReplace   = fmt.Sprintf("%v must be replaced", resourceName)
)

// Handle line that indicates the start of a Terraform plan:
// "Terraform will perform the following actions:"
func parseStartPlan(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(resourceCreated, Line)
	if !match {
		return false, nil
	}
	log.ContainsPlan = true
	return true, nil
}

// Handle line that indicates a resource is tainted. E.g:
// "  # aws_ssm_parameter.p1 is tainted, so must be replaced"
func parsePlanTainted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(isTainted, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse tainted resource: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " is tainted, so must be replaced")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, Created)
	return true, nil
}

// Handle line that indicates a resource has been marked to be replaced. E.g:
// "  # aws_ssm_parameter.p1 will be replaced, as requested"
func parsePlanExplicitReplace(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(explicitReplace, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource to replace: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " will be replaced, as requested")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, Created)
	return true, nil
}

// Handle line that indicates a resource will be destroyed. E.g:
// "  # aws_ssm_parameter.p1 will be destroyed"
func parsePlanWillBeDestroyed(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(willBeDestroyed, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource for destroy: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " will be destroyed")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, NotCreated)
	return true, nil
}

// Handle line that indicates a resource will be modified. E.g:
// " # aws_ssm_parameter.p5 will be updated in-place"
func parsePlanWillBeModified(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(willBeModified, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource for modify: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " will be updated in-place")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, Created)
	return true, nil
}

// Handle line that indicates a resource must be replaced. E.g:
// "# aws_ssm_parameter.p6 must be replaced"
func parsePlanForcedReplace(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(forcedReplace, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse resource for replacement: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " must be replaced")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, Created)
	return true, nil
}

// Handle line that indicates a resource will be created. E.g:
// "# aws_ssm_parameter.p6 will be created"
func parsePlanWillBeCreated(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(willBeCreated, Line)
	if !match {
		return false, nil
	}

	tokens := strings.Split(Line, "# ")
	if len(tokens) < 2 {
		msg := fmt.Sprintf("Unable to parse creation plan: %v\n", Line)
		return false, &LineParseError{Msg: msg}
	}
	resource := strings.Split(tokens[1], " will be created")[0]

	log.RegisterNewResource(resource)
	log.SetDesiredStatus(resource, Created)
	return true, nil
}
