package mold

import (
	"log"

	amold "github.com/ArthurKC/foundry/pkg/adapters/mold"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init TARGET_DIR",
	Args:  cobra.ExactArgs(1),
	Short: "Init directory as mold.",
	Long:  `Init directory as mold. Create mold.yaml meta file.`,
	Run: func(cmd *cobra.Command, args []string) {
		meta := amold.NewMetaFile(args[0])
		if err := meta.Initialize(); err != nil {
			log.Fatal(err)
		}
	},
}
