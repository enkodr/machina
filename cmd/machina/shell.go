package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var shellCommand = &cobra.Command{
	Use:     "shell",
	Short:   "Enters a instance shell",
	Aliases: []string{"sh"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the instance data
		instance, err := hypvsr.GetMachine(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Instance %q does not exist\n", name)
			os.Exit(1)
		}

		// Entering the instance shell
		fmt.Println(fmt.Sprintf("Enter instance's '%s' shell\n", name))
		err = instance.Shell()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error entering instance's %q shell\n", name)
			os.Exit(1)
		}
	},
}

func init() {
	rootCommand.AddCommand(shellCommand)
}
