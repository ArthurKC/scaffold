package template

import (
	"log"

	atemplate "github.com/ArthurKC/scaffold/pkg/adapters/template"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init TARGET_DIR",
	Args:  cobra.ExactArgs(1),
	Short: "Init directory as template.",
	Long:  `Init directory as template. Create scaffold.yaml meta file.`,
	Run: func(cmd *cobra.Command, args []string) {
		meta := atemplate.NewMetaFile(args[0])
		if err := meta.Initialize(); err != nil {
			log.Fatal(err)
		}
	},
}
