package machina

import (
	"fmt"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var deleteCommand = &cobra.Command{
	Use:   "delete",
	Short: "Delete a machine",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		vm, err := vm.Load(name)
		if err != nil {
			log.Error("the machine doesn't exist")
		}

		// Start the machine
		log.Info(fmt.Sprintf("Deleting machine '%s'", name))
		vm.Delete()

	},
}

func init() {
	rootCommand.AddCommand(deleteCommand)
}
