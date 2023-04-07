package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(graphCmd)
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Visualize a Terraform run graphically",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("'graph' command has not been implemented yet!")
	},
}
