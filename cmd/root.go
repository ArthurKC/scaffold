package cmd

import (
	"fmt"
	"os"

	"github.com/ArthurKC/scaffold/cmd/template"

	"github.com/spf13/cobra"
)

func init() {
	template.BindCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "generate scaffold.",
	Long:  `generate scaffold.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
