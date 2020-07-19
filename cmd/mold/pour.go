package mold

import (
	"os"
	"strings"

	"github.com/ArthurKC/foundry/pkg/adapters/material"
	amold "github.com/ArthurKC/foundry/pkg/adapters/mold"
	umold "github.com/ArthurKC/foundry/pkg/usecases/mold"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pourCmd)
}

var pourCmd = &cobra.Command{
	Use:   "pour MOLD_DIR DEST_DIR",
	Args:  cobra.ExactArgs(2),
	Short: "pour molten material into a mold.",
	Long:  `pour molten material into a mold.`,
	Run: func(cmd *cobra.Command, args []string) {
		moldDir := strings.TrimSuffix(args[0], "/")
		destDir := strings.TrimSuffix(args[1], "/")
		pour := umold.NewPourInteractor(
			amold.NewFileRepository(),
			material.NewIOMaterialService(os.Stdout, os.Stdin),
			amold.NewFileProductService(os.Stdout),
			amold.NewIOOutputPort(os.Stdout),
		)
		pour.ExecutePour(&umold.PourRequest{
			MoldName: moldDir,
			DestDir:  destDir,
		})
	},
}
