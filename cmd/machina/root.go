package machina

import (
	"github.com/spf13/cobra"
)

var (
	name string
	file string
)

var rootCommand = &cobra.Command{
	Use:   "machina",
	Short: "A tool to manage cloud images with KVM",
}

func Execute() {
	rootCommand.Execute()
}

func validateName(cmd *cobra.Command, args []string) error {
	name = "default"
	return nil
}
