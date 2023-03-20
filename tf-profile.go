package main

import (
	"log"
	"os"

	cli "github.com/QuintenBruynseraede/tf-profile/cli"
)

// Main entrypoint to the CLI
func main() {
	tfprofile := cli.Create()
	err := tfprofile.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
