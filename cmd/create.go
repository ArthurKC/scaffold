package cmd

import (
	agenerator "github.com/ArthurKC/scaffold/pkg/adapters/generator"
	"github.com/ArthurKC/scaffold/pkg/domains/generator"
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
		tSrc, err := agenerator.NewTemplateSource(templateDir)
		if err != nil {
			panic(err)
		}
		g := generator.New(tSrc, agenerator.NewInput(), agenerator.NewOutput(destDir))
		g.Generate()
	},
}
