package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "generate scaffold.",
	Long:  `generate scaffold.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello scaffold")
	},
}

func main() {
	rootCmd.Execute()
}
