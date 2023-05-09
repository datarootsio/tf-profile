package cmd

import (
	stats "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/stats"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().BoolP("tee", "t", false, "Print logs while parsing")
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Parse a Terraform log and show general statistics",
	Long: `The 'stats' command can be used to show general statistics 
	a Terraform run. It prints high-level statistics on the following topics:
	basic, time-related, creation status and modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return stats.Stats(args, tee)
	},
}
