package cmd

import (
	"fmt"
	"os"

	"github.com/ArthurKC/foundry/cmd/material"
	"github.com/ArthurKC/foundry/cmd/mold"

	"github.com/spf13/cobra"
)

func init() {
	material.BindCommand(rootCmd)
	mold.BindCommand(rootCmd)
}

var rootCmd = &cobra.Command{
	Use:   "foundry",
	Short: "generate scaffold.",
	Long:  `generate scaffold.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
