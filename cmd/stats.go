package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Parse a Terraform log and show general statistics",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("'stats' command has not been implemented yet!")
	},
}
