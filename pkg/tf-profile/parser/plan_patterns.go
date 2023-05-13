package tfprofile

import (
	"regexp"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

var (
	StartPlan = "Terraform will perform the following actions:"
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
