package mold

import (
	"os"
	"strings"

	amold "github.com/ArthurKC/foundry/pkg/adapters/mold"
	umold "github.com/ArthurKC/foundry/pkg/usecases/mold"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "create MOLD_DIR [FROM_DIR]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Create directory as mold.",
	Long:  `Create directory as mold. Create mold.yaml meta file.`,
	Run: func(cmd *cobra.Command, args []string) {
		moldDir := strings.TrimSuffix(args[0], "/")
		importDir := moldDir
		if len(args) >= 2 {
			importDir = strings.TrimSuffix(args[1], "/")
		}
		u := umold.NewCreateInteractor(
			amold.NewFileRepository(),
			amold.NewFileService(),
			amold.NewIOOutputPort(os.Stdout),
		)
		u.ExecuteCreate(&umold.CreateRequest{
			ImportPath: importDir,
			MoldName:   moldDir,
		})
	},
}
