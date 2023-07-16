package tfprofile

import "regexp"

// Terraform inserts a lot of formatting strings into its output when
// -no-color is not specified. This function removes all of those
func RemoveTerminalFormatting(in string) string {
	// regex to detect ANSI terminal formatting directives (https://stackoverflow.com/a/14693789)
	re := regexp.MustCompile(`(?:\x1B[@-_]|[\x80-\x9F])[0-?]*[ -/]*[@-~]`)
	return re.ReplaceAllString(in, "")
}
