package mold

import (
	"github.com/spf13/cobra"
)

func BindCommand(parent *cobra.Command) {
	parent.AddCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "mold",
	Short: "Manage mold.",
	Long:  `Manage mold.`,
}
