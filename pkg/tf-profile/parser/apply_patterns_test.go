package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestParseCreate(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := parseResourceCreationStarted("foo: Creating...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = parseResourceCreated("foo: Creation complete after 1s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(1000), log.Resources["foo"].TotalTime)
	assert.Equal(t, Created, log.Resources["foo"].AfterStatus)
}

func TestParseCreateFailed(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := parseResourceCreationStarted("foo: Creating...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = parseResourceCreationFailed("with foo,", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Failed, log.Resources["foo"].AfterStatus)
}

func TestResourceDestruction(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := parseResourceDestructionStarted("foo: Destroying...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = parseResourceDestroyed("foo: Destruction complete after 10s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(10000), log.Resources["foo"].TotalTime)
	assert.Equal(t, NotCreated, log.Resources["foo"].AfterStatus)
}

func TestResourceModification(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := parseResourceModificationStarted("foo: Modifying...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = parseResourceModified("foo: Modifications complete after 10s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(10000), log.Resources["foo"].TotalTime)
	assert.Equal(t, Created, log.Resources["foo"].AfterStatus)
}
