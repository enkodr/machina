package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var (
	force bool
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running machine",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		machine, err := hypvsr.Load(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Machine %q doesn't exist\n", name)
			os.Exit(1)
		}

		// Start the machine
		fmt.Printf("Stoping machine '%s'\n", name)
		if force {
			err = machine.ForceStop()
		} else {
			err = machine.Stop()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error stoping the machine information\n")
			os.Exit(1)
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	stopCommand.PersistentFlags().BoolVarP(&force, "force", "f", false, "force stop")
	rootCommand.AddCommand(stopCommand)
}
