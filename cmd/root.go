package cmd

import (
	"fmt"
	"os"

	"github.com/ArthurKC/foundry/cmd/metal"
	"github.com/ArthurKC/foundry/cmd/mold"

	"github.com/spf13/cobra"
)

func init() {
	metal.BindCommand(rootCmd)
	mold.BindCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "foundry",
	Short: "generate foundry.",
	Long:  `generate foundry.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
