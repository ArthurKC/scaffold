package template

import (
	"github.com/spf13/cobra"
)

func BindCommand(parent *cobra.Command) {
	parent.AddCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage template.",
	Long:  `Manage template.`,
}
