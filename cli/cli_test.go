package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppCreate(t *testing.T) {
	app := Create()

	// Only basic checks on the app, we don't retest a dependency
	assert.Equal(t, "tf-profile", app.Name)
	assert.Equal(t, 5, len(app.Flags))
}

func TestCorrectArguments(t *testing.T) {
	app := Create()
	assert.Nil(t, app.Run([]string{"../tf-profile", "--tee"}))
	assert.Nil(t, app.Run([]string{"../tf-profile"}))
	assert.Nil(t, app.Run([]string{"../tf-profile", "../test_files/multiple_resources.log"}))
}

func TestIncorrectArguments(t *testing.T) {
	app := Create()
	err1 := app.Run([]string{"./tf-profile", "--stats"})
	assert.NotNil(t, err1) // --stats is not implemented yet

	err2 := app.Run([]string{"./tf-profile", "--max_depth=1234"})
	assert.NotNil(t, err2) // --max_depth is not implemented yet

	err3 := app.Run([]string{"./tf-profile", "arg1", "arg2"})
	assert.NotNil(t, err3) // Only one input file is supported

}

func TestBasicRun(t *testing.T) {
	args := &InputArgs{
		debug:      true,
		log_level:  "INFO",
		stats:      false,
		tee:        true,
		max_depth:  1,
		sort:       "tot_time=asc",
		input_file: "", // stdin
	}

	printArgs(args)
	err := run(args)
	assert.Nil(t, err)
}

func TestFileDoesntExist(t *testing.T) {
	args := &InputArgs{
		debug:      true,
		log_level:  "INFO",
		stats:      false,
		tee:        true,
		max_depth:  1,
		sort:       "tot_time=asc",
		input_file: "does-not-exist",
	}
	err := run(args)
	assert.NotNil(t, err)
}
