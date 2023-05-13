package tfprofile

type (
	LineParseError struct{ Msg string }
)

func (e *LineParseError) Error() string {
	return e.Msg
}
