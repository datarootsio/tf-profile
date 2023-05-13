package tfprofile

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

var (
	StartPlan       = "Terraform will perform the following actions:"
	IsTainted       = fmt.Sprintf("%v is tainted, so must be replaced", ResourceName)
	WillBeCreated   = fmt.Sprintf("%v will be created", ResourceName)
	ExplicitReplace = fmt.Sprintf("%v will be replaced, as requested", ResourceName)
	WillBeDestroyed = fmt.Sprintf("%v will be destroyed", ResourceName)
	WillBeModified  = fmt.Sprintf("%v will be updated in-place", ResourceName)
	ForcedReplace   = fmt.Sprintf("%v must be replaced", ResourceName)
)

// Handle line that indicates the start of a Terraform plan:
// "Terraform will perform the following actions:"
func ParseStartPlan(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ResourceCreated, Line)
	if !match {
		return false, nil
	}
	log.ContainsPlan = true
	return true, nil
}

// Handle line that indicates a resource is tainted. E.g:
// "  # aws_ssm_parameter.p1 is tainted, so must be replaced"
func ParsePlanTainted(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(IsTainted, Line)
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
func ParsePlanExplicitReplace(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ExplicitReplace, Line)
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
func ParsePlanWillBeDestroyed(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(WillBeDestroyed, Line)
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
func ParsePlanWillBeModified(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(WillBeModified, Line)
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
func ParsePlanForcedReplace(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(ForcedReplace, Line)
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
func ParsePlanWillBeCreated(Line string, log *ParsedLog) (bool, error) {
	match, _ := regexp.MatchString(WillBeCreated, Line)
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
