package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit.",
	Long:  `Print the version and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tf-profile v0.5.0")
	},
}
