package metal

import (
	ametal "github.com/ArthurKC/foundry/pkg/adapters/metal"
	"github.com/ArthurKC/foundry/pkg/domains/metal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pourCmd)
}

var pourCmd = &cobra.Command{
	Use:   "pour_into MOLD_DIR DEST_DIR",
	Args:  cobra.ExactArgs(2),
	Short: "pour molten metal into a mold.",
	Long:  `pour molten metal into a mold.`,
	Run: func(cmd *cobra.Command, args []string) {
		moldDir := args[0]
		destDir := args[1]
		tSrc, err := ametal.NewMoldSource(moldDir)
		if err != nil {
			panic(err)
		}
		g := metal.New(tSrc, ametal.NewInput(), ametal.NewOutput(destDir))
		g.Generate()
	},
}
