package main

// // Main entrypoint to the CLI
// func main() {
// 	tfprofile := tfprofile.Create()
// 	err := tfprofile.Run(os.Args)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

import (
	"github.com/QuintenBruynseraede/tf-profile/cmd"
)

func main() {
	cmd.Execute()
}
