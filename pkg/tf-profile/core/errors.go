package tfprofile

import "fmt"

type (
	LineParseError        struct{ Msg string }
	ResourceNotFoundError struct{ Resource string }
)

func (e *LineParseError) Error() string {
	return e.Msg
}

func (e *ResourceNotFoundError) Error() string {
	return fmt.Sprintf("Unable to find resource %v in log.", e.Resource)
}
