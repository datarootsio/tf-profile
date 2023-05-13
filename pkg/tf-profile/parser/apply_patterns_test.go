package tfprofile

import (
	"testing"

	. "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/core"
	"github.com/stretchr/testify/assert"
)

func TestParseCreate(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := ParseResourceCreationStarted("foo: Creating...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = ParseResourceCreated("foo: Creation complete after 1s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(1000), log.Resources["foo"].TotalTime)
	assert.Equal(t, Created, log.Resources["foo"].AfterStatus)
}

func TestParseCreateFailed(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := ParseResourceCreationStarted("foo: Creating...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = ParseResourceCreationFailed("with foo,", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, Failed, log.Resources["foo"].AfterStatus)
}

func TestResourceDestruction(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := ParseResourceDestructionStarted("foo: Destroying...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = ParseResourceDestroyed("foo: Destruction complete after 10s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(10000), log.Resources["foo"].TotalTime)
	assert.Equal(t, NotCreated, log.Resources["foo"].AfterStatus)
}

func TestResourceModification(t *testing.T) {
	log := ParsedLog{Resources: map[string]ResourceMetric{}}

	modified, err := ParseResourceModificationStarted("foo: Modifying...", &log)
	assert.True(t, modified)
	assert.Nil(t, err)

	modified, err = ParseResourceModified("foo: Modifications complete after 10s [id=/no/slash/at/end0]", &log)
	assert.True(t, modified)
	assert.Nil(t, err)
	assert.Equal(t, float64(10000), log.Resources["foo"].TotalTime)
	assert.Equal(t, Created, log.Resources["foo"].AfterStatus)
}
