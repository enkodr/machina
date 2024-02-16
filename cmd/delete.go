package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var (
	forceDelete bool
)

var deleteCommand = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a instance",
	Aliases: []string{"rm"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {

		// Get the instance name from the first argument
		for _, arg := range args {
			name = arg
			// Load the instance data
			instance, err := hypvsr.GetMachine(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "the instance %q doesn't exist\n", name)
				return
			}

			// Confirm if the instance will be deleted
			if !forceDelete {
				// Ask the user for delete confirmation
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("Are you certain you want to delete instance '%s' [y/N]: ", name)
				response, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
					os.Exit(1)
				}

				// Convert the user response to lowercase
				response = strings.ToLower(strings.TrimSpace(response))

				if response == "y" || response == "yes" {
					forceDelete = true
				}

			}

			// Check's if the flag to force delete was passed
			// or if the users confirmed the deletion action
			if forceDelete {
				// Delete the instance
				fmt.Printf("Deleting instance '%s'\n", name)
				err = instance.Delete()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting the instance\n")
				}
			}
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	deleteCommand.PersistentFlags().BoolVarP(&forceDelete, "yes", "y", false, "yes will be assumed and you wont recieve a confirmation prompt")
	rootCommand.AddCommand(deleteCommand)
}
