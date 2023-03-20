package cli

import (
	"bufio"
	"errors"
	"fmt"
	"log"

	"github.com/QuintenBruynseraede/tf-profile/parser"
	print "github.com/QuintenBruynseraede/tf-profile/printer"
	"github.com/QuintenBruynseraede/tf-profile/reader"

	"github.com/urfave/cli"
)

func Create() *cli.App {
	return &cli.App{
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
			args, err := parseArgs(c)
			if err != nil {
				log.Fatalf("Error during argument parsing: \n%v\n", err)
			}

			err2 := validateArgs(args)
			if err2 != nil {
				log.Fatalf("Error during argument validation: \n%v\n", err)
			}

			if args.debug {
				printArgs(args)
			}

			err3 := run(args)
			if err3 != nil {
				log.Fatalf("Error during tf-profile run:\n%v\n", err)
			}

			return nil
		},
	}
}

type InputArgs struct {
	debug      bool
	log_level  string
	stats      bool
	tee        bool
	max_depth  int
	sort       string
	input_file string
}

func parseArgs(c *cli.Context) (*InputArgs, error) {
	var input_file string
	if c.NArg() == 1 {
		input_file = c.Args().Get(0)
	} else {
		input_file = ""
	}

	if c.NArg() > 1 {
		msg := fmt.Sprintf("Expected at most 1 argument, received %v: %v\n", c.NArg(), c.Args())
		return nil, errors.New(msg)
	}

	return &InputArgs{
		debug:      c.Bool("debug"),
		log_level:  c.String("log_level"),
		stats:      c.Bool("stats"),
		tee:        c.Bool("tee"),
		max_depth:  c.Int("max_depth"),
		sort:       c.String("sort"),
		input_file: input_file,
	}, nil
}

func printArgs(args *InputArgs) {
	fmt.Println("==== tf-profile ====")
	fmt.Printf("Running with config:\n")
	fmt.Printf("- log_level: %v\n", args.log_level)
	fmt.Printf("- stats: %v\n", args.stats)
	fmt.Printf("- tee: %v\n", args.tee)
	fmt.Printf("- max_depth: %v\n", args.max_depth)
	fmt.Printf("- sort: %v\n", args.sort)
	fmt.Println("====================")
}

// Validate all arguments passed into the CLI tool
// will print an error message and exit with a non-zero
// exitcode if incompatible arguments are detected.
func validateArgs(args *InputArgs) error {
	if args.max_depth != -1 {
		log.Fatal("--max_depth is not implemented yet!")
	}
	if args.stats {
		log.Fatal("--stats is not implemented yet!")
	}

	// TODO: check that the file comes last, i.e. tf-profile --tee logs.txt | NOT tf-profile logs.txt --tee
	// TODO: check spec format
	return nil
}

func run(args *InputArgs) error {
	var file *bufio.Scanner

	if args.input_file != "" {
		fmt.Printf("Input: from file %v\n", args.input_file)
		file = reader.FileReader{File: args.input_file}.Read()
	} else {
		fmt.Printf("Input: from stdin\n")
		file = reader.StdinReader{}.Read()
	}

	tflog := parser.Parse(file, args.tee)

	print.Table(&tflog, args.sort)
	return nil
}
