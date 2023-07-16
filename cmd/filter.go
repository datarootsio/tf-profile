package cmd

import (
	filter "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/filter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(filterCmd)
}

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter a Terraform log to only selected resources",
	Long: `The 'filter' command is used to filter down logs to only
those lines that contain references to a set of selected resources.
Resources can be specified with regex. Only lines containing those
resources will be printed. Terraform plans pertaining this resource
are always fully shown.

This command expects two arguments when reading from a log file:

$ tf-profile filter "aws_ssm_parameter.*" /path/to/log.txt

Only one argument is needed to read from stdin:

$ terraform apply -auto-approve | tf-profile filter "aws_ssm_parameter.*"  

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return filter.Filter(args)
	},
}
