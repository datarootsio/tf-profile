package tfprofile

import (
	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
)

// Doesn't parse anything
func DummyParser(Line string, log *ParsedLog) (bool, error) {
	return true, nil
}
