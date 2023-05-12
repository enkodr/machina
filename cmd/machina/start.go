package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/vm"
	"github.com/spf13/cobra"
)

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts a stopped machine",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		machine, err := vm.Load(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "the machine %q doesn't exist\n", name)
			os.Exit(1)
		}

		// Start the machine
		fmt.Printf("Starting machine '%s'\n", name)
		err = machine.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create the machine\n")
			os.Exit(1)
		}

		// Wait until machine is ready
		fmt.Printf("Waiting for machine '%s'\n", name)
		err = machine.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Machine seems to be stuck in start process\n")
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	rootCommand.AddCommand(startCommand)
}
