package cmd

import (
	table "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/table"
	"github.com/spf13/cobra"
)

var (
	max_depth int
	tee       bool
	sort      string
)

func init() {
	rootCmd.AddCommand(tableCmd)
	tableCmd.Flags().StringVarP(
		&sort,
		"sort",
		"s",
		"tot_time=desc,resource=asc",
		"Comma-separated list of KEY=(asc|desc) to control sorting.",
	)
	tableCmd.Flags().IntVarP(
		&max_depth,
		"max_depth",
		"d",
		-1,
		"Max recursive module depth before aggregating.",
	)
	tableCmd.Flags().Bool("tee", false, "Print logs while parsing")
}

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "Parse a Terraform log and perform resource-level profiling",
	Args:  cobra.MaximumNArgs(1),
	Long: `The 'table' command is used to do in-depth profiling on a resource level.
	It will parse a log, extract metrics about all resources and show tabular output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return table.Table(args, max_depth, tee, sort)
	},
}
