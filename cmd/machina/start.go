package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var startCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts a stopped instance",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			name = arg
			// Load the instance data
			instance, err := hypvsr.GetMachine(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "the instance %q doesn't exist\n", name)
				os.Exit(1)
			}

			// Start the instance
			fmt.Printf("Starting instance '%s'\n", name)
			err = instance.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to create the instance\n")
				os.Exit(1)
			}

			// Wait until instance is ready
			fmt.Printf("Waiting for instance '%s'\n", name)
			err = instance.Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Instance seems to be stuck in start process\n")
			}

		}
		fmt.Printf("Done!\n")
	},
}

func init() {
	rootCommand.AddCommand(startCommand)
}
