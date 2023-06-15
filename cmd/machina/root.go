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

// Validate instance name and define the default name
func validateName(cmd *cobra.Command, args []string) error {
	name = "default"
	if len(args) > 0 {
		name = args[0]
	}
	return nil
}
