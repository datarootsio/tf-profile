package tfprofile

import (
	"bufio"
)

// Execute the `tf-profile table` command
func Table(args []string, max_depth int, tee bool, sort string) error {
	var file *bufio.Scanner
	var err error

	if len(args) == 1 {
		file, err = FileReader{File: args[0]}.Read()
	} else {
		file, err = StdinReader{}.Read()
	}

	if err != nil {
		return err
	}

	tflog, err := Parse(file, tee)
	if err != nil {
		return err
	}

	tflog, err = Aggregate(tflog)
	if err != nil {
		return err
	}

	err = PrintTable(tflog, sort)
	if err != nil {
		return err
	}

	return nil
}
