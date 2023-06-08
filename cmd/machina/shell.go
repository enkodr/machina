package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var shellCommand = &cobra.Command{
	Use:     "shell",
	Short:   "Enters a machine shell",
	Aliases: []string{"sh"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		machine, err := hypvsr.Load(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Machine %q does not exist\n", name)
			os.Exit(1)
		}

		// Entering the machine shell
		fmt.Println(fmt.Sprintf("Enter machine's '%s' shell\n", name))
		err = machine.Shell()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error entering machine's %q shell\n", name)
			os.Exit(1)
		}
	},
}

func init() {
	rootCommand.AddCommand(shellCommand)
}
