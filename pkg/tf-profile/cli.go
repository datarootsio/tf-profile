package tfprofile

import (
	"bufio"
	"errors"
	"fmt"

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
				return errors.New("Error during argument parsing")
			}

			err2 := validateArgs(args)
			if err2 != nil {
				return errors.New("Error during argument validation")
			}

			if args.debug {
				printArgs(args)
			}

			err3 := run(args)
			if err3 != nil {
				return errors.New("Error during tf-profile run")
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
		return errors.New("--max_depth is not implemented yet!")
	}
	if args.stats {
		return errors.New("--stats is not implemented yet!")
	}

	// TODO: check that the file comes last, i.e. tf-profile --tee logs.txt | NOT tf-profile logs.txt --tee
	// TODO: check spec format
	return nil
}

func run(args *InputArgs) error {
	var file *bufio.Scanner
	var err1 error

	if args.input_file != "" {
		if args.debug {
			fmt.Printf("Input: from file %v\n", args.input_file)
		}
		file, err1 = FileReader{File: args.input_file}.Read()
	} else {
		if args.debug {
			fmt.Printf("Input: from stdin\n")

		}
		file, err1 = StdinReader{}.Read()
	}

	if err1 != nil {
		return err1
	}

	tflog, err2 := Parse(file, args.tee)
	if err2 != nil {
		return err2
	}

	err3 := Table(&tflog, args.sort)
	if err3 != nil {
		return err3
	}

	return nil
}
