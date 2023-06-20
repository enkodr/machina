package machina

import (
	"fmt"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var copyCommand = &cobra.Command{
	Use:               "copy",
	Short:             "Copies files/directories from the host to the VM and vice-versa",
	Aliases:           []string{"cp"},
	ValidArgsFunction: bashCompleteInstanceNamesConnection,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the instance data
		if len(args) > 0 {
			name = args[0]
		}

		// Check if all required arguments are passed
		if len(args) != 2 {
			fmt.Printf("machina copy <host_path> <instance_name>:<instance_path>\n")
			fmt.Printf("machina copy <instance_name>:<instance_path> <host_path>\n")
			os.Exit(0)
		}

		// As the user needs to pass two arguments and those can be in any direction
		// (host -> VM or VM -> host), this logic will identify the direction
		origin := args[0]
		dest := args[1]

		// Identify the copy direction based on the existence
		// on the existence of a colon
		if strings.Contains(origin, ":") {
			name = strings.Split(origin, ":")[0]
		} else {
			name = strings.Split(dest, ":")[0]
		}

		// Get the instance from the configuration file
		instance, err := hypvsr.GetMachine(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Instance %q does not exist\n", name)
			os.Exit(1)
		}

		// Copy the content from origin to destination
		fmt.Printf("Copying content...\n")
		err = instance.CopyContent(origin, dest)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error copying content\n")
			os.Exit(1)
		}

		fmt.Printf("Done!\n")

	},
}

func init() {
	rootCommand.AddCommand(copyCommand)
}
