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
	Short: "Stops a running instance",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the instance data
		instance, err := hypvsr.GetMachine(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Instance %q doesn't exist\n", name)
			os.Exit(1)
		}

		// Stop the instance
		fmt.Printf("Stoping instance '%s'\n", name)
		if force {
			err = instance.ForceStop()
		} else {
			err = instance.Stop()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error stoping the instance information\n")
			os.Exit(1)
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	stopCommand.PersistentFlags().BoolVarP(&force, "force", "f", false, "force stop")
	rootCommand.AddCommand(stopCommand)
}
