package material

import (
	"github.com/spf13/cobra"
)

func BindCommand(parent *cobra.Command) {
	parent.AddCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "material",
	Short: "Manage material.",
	Long:  `Manage material.`,
}
