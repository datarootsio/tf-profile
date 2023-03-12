package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/QuintenBruynseraede/tf-profile/parser"
	print "github.com/QuintenBruynseraede/tf-profile/printer"
	"github.com/QuintenBruynseraede/tf-profile/reader"

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
			cli.StringFlag{
				Name:  "sort",
				Usage: "Sort specification",
				Value: "tot_time=desc,idx_created=asc",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println("==== tf-profile ====")
			fmt.Printf("Running with config:\n")
			fmt.Printf("Received %v args, %v flags\n", c.NArg(), c.NumFlags())
			fmt.Printf("- log_level: %v\n", c.String("log_level"))
			fmt.Printf("- stats: %v\n", c.Bool("stats"))
			fmt.Printf("- tee: %v\n", c.Bool("tee"))
			fmt.Printf("- max_depth: %v\n", c.Int("max_depth"))
			fmt.Printf("- sort: %v\n", c.String("sort"))
			fmt.Println("====================")

			ValidateArgs(c)
			Run(c)

			return nil
		},
	}

	err := tfprofile.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Validate all arguments passed into the CLI tool
// will print an error message and exit with a non-zero
// exitcode if incompatible arguments are detected.
func ValidateArgs(c *cli.Context) {
	// TODO: check that the file comes last, i.e. tf-profile --tee logs.txt | NOT tf-profile logs.txt --tee
	// TODO: sort spec format
}

func Run(c *cli.Context) {
	inputFile := ""
	var file *bufio.Scanner

	c.NumFlags()
	if c.NArg() == 1 {
		inputFile = c.Args().Get(0)
		fmt.Printf("Input: from file %v\n", inputFile)
		file = reader.FileReader{File: inputFile}.Read()
	} else {
		fmt.Printf("Input: from stdin\n")
		file = reader.StdinReader{}.Read()
	}

	tflog := parser.Parse(file, c.Bool("tee"))

	print.Table(&tflog, c.String("sort"))
}
