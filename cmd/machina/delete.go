package machina

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/vm"
	"github.com/spf13/cobra"
)

var (
	forceDelete bool
)

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a machine",
	Aliases: []string{"rm"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			name = args[0]
		}
		// Load the machine data
		vm, err := vm.Load(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "the machine %q doesn't exist\n", name)
			return
		}

		// Confirm if the machine will be deleted
		if !forceDelete {
			// Ask the user for delete confirmation
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Are you certain you want to delete machine '%s' [y/N]: ", name)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}

			response = strings.ToLower(strings.TrimSpace(response))

			if response == "y" || response == "yes" {
				forceDelete = true
			}

		}
		if forceDelete {
			// Delete the machine
			fmt.Printf("Deleting machine '%s'\n", name)
			err = vm.Delete()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting the machine\n")
			}
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	deleteCommand.PersistentFlags().BoolVarP(&forceDelete, "yes", "y", false, "yes will be assumed and you wont recieve a confirmation prompt")
	rootCommand.AddCommand(deleteCommand)
}
