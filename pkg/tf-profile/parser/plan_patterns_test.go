package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestParsePlan(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := ParsePlanTainted("  # foo is tainted, so must be replaced", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = ParsePlanExplicitReplace("  # foo will be replaced, as requested", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = ParsePlanWillBeDestroyed("  # foo will be destroyed", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, NotCreated, log.Resources["foo"].DesiredStatus)

	modified, err = ParsePlanWillBeModified("  # foo will be updated in-place", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = ParsePlanForcedReplace("  # foo must be replaced", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)

	modified, err = ParsePlanWillBeCreated("  # foo will be created", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Created, log.Resources["foo"].DesiredStatus)
}

// Not the best test as we need to construct text that passes the regex check,
// but still throws an error during parsing
func TestParseErrors(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}
	_, err := ParsePlanTainted("foo is tainted, so must be replaced", &log)
	assert.NotNil(t, err)
	_, err = ParsePlanExplicitReplace("foo will be replaced, as requested", &log)
	assert.NotNil(t, err)
	_, err = ParsePlanWillBeDestroyed("foo will be destroyed", &log)
	assert.NotNil(t, err)
	_, err = ParsePlanWillBeModified("foo will be updated in-place", &log)
	assert.NotNil(t, err)
	_, err = ParsePlanForcedReplace("foo must be replaced", &log)
	assert.NotNil(t, err)
	_, err = ParsePlanWillBeCreated("foo will be created", &log)
	assert.NotNil(t, err)
}
