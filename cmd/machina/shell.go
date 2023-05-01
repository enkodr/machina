package machina

import (
	"fmt"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var shellCommand = &cobra.Command{
	Use:   "shell",
	Short: "Enters a machine shell",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		machine, err := vm.Load(name)
		if err != nil {
			log.Error("the machine doesn't exist")
		}

		// Start the machine
		log.Info(fmt.Sprintf("Enter the machine '%s' shell", name))
		machine.Shell()

	},
}

func init() {
	rootCommand.AddCommand(shellCommand)
}
