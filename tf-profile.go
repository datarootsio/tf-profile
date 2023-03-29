package main

import (
	"log"
	"os"

	tfprofile "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile"
)

// Main entrypoint to the CLI
func main() {
	tfprofile := tfprofile.Create()
	err := tfprofile.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
