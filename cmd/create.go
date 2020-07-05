package cmd

import (
	"github.com/ArthurKC/scaffold/pkg/cui"
	"github.com/ArthurKC/scaffold/pkg/file"
	"github.com/ArthurKC/scaffold/pkg/generator"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create TEMPLATE_DIR DEST_DIR",
	Args:  cobra.ExactArgs(2),
	Short: "generate scaffold.",
	Long:  `generate scaffold.`,
	Run: func(cmd *cobra.Command, args []string) {
		templateDir := args[0]
		destDir := args[1]
		tSrc, err := file.NewTemplateSource(templateDir)
		if err != nil {
			panic(err)
		}
		g := generator.New(tSrc, cui.NewInput(), file.NewOutput(destDir))
		g.Generate()
	},
}
