package main

import (
	"fmt"
	"log"
	"os"

	"github.com/QuintenBruynseraede/tf-profile/readers"
	"github.com/urfave/cli"
)

// Main entrypoint to the CLI
func main() {
	var tfprofile = cli.App{
		Name:    "tf-profile",
		Usage:   "CLI tool to profile Terraform runs, written in Go",
		Author:  "Quinten Bruynseraede",
		Version: "0.0.1",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "log_level",
				Value:  "INFO",
				Usage:  "cli log level as read from TF_LOG",
				EnvVar: "TF_LOG",
			},
			cli.BoolFlag{
				Name:  "stats",
				Usage: "Show global stats only",
			},
			cli.IntFlag{
				Name:  "max_depth",
				Value: -1,
				Usage: "Max depth of submodules before aggregating metrics.",
			},
			cli.BoolFlag{
				Name:  "tee",
				Usage: "print to stdout while profiling",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("==== tf-profile ====")
			fmt.Printf("Running with config:\n")
			fmt.Printf("- log_level: %v\n", c.String("log_level"))
			fmt.Printf("- stats: %v\n", c.Bool("stats"))
			fmt.Printf("- tee: %v\n", c.Bool("tee"))
			fmt.Printf("- max_depth: %v\n", c.Int("max_depth"))
			fmt.Println("====================")

			inputFile := ""
			if c.NArg() == 1 {
				inputFile = c.Args().Get(0)
				fmt.Printf("Input: from file %v\n", inputFile)
			} else {
				fmt.Printf("Input: from stdin\n")
				reader := readers.StdinReader{Tee: true}
				reader.ReadFile()
			}

			return nil
		},
	}

	err := tfprofile.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
