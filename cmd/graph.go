package cmd

import (
	"fmt"

	graph "github.com/QuintenBruynseraede/tf-profile/pkg/tf-profile/graph"

	"github.com/spf13/cobra"
)

var (
	Size    []int
	OutFile string
)

func init() {
	rootCmd.AddCommand(graphCmd)
	graphCmd.Flags().IntSliceVarP(&Size, "size", "s", []int{1000, 600}, "Width and height of generated image")
	graphCmd.Flags().StringVarP(&OutFile, "out", "o", "tf-profile-graph.png", "Output file used by gnuplot")
	graphCmd.Flags().BoolVarP(&aggregate, "aggregate", "a", true, "Agregate count[] and for_each[]")
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Visualize a Terraform run graphically",
	Long:  `TODO`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(Size) != 2 || Size[0] < 0 || Size[1] < 0 {
			return fmt.Errorf("Expected two positive integers for --size flag, got %v", Size)
		}
		return graph.Graph(args, Size[0], Size[1], OutFile, aggregate)
	},
}
