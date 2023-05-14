package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestParsePlan(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := parsePlanTainted("  # foo is tainted, so must be replaced", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = parsePlanExplicitReplace("  # foo will be replaced, as requested", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = parsePlanWillBeDestroyed("  # foo will be destroyed", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, NotCreated, log.Resources["foo"].DesiredStatus)

	modified, err = parsePlanWillBeModified("  # foo will be updated in-place", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = parsePlanForcedReplace("  # foo must be replaced", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = parsePlanWillBeCreated("  # foo will be created", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)
}

// Not the best test as we need to construct text that passes the regex check,
// but still throws an error during parsing
func TestParseErrors(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}
	_, err := parsePlanTainted("foo is tainted, so must be replaced", &log)
	assert.NotNil(t, err)
	_, err = parsePlanExplicitReplace("foo will be replaced, as requested", &log)
	assert.NotNil(t, err)
	_, err = parsePlanWillBeDestroyed("foo will be destroyed", &log)
	assert.NotNil(t, err)
	_, err = parsePlanWillBeModified("foo will be updated in-place", &log)
	assert.NotNil(t, err)
	_, err = parsePlanForcedReplace("foo must be replaced", &log)
	assert.NotNil(t, err)
	_, err = parsePlanWillBeCreated("foo will be created", &log)
	assert.NotNil(t, err)
}
