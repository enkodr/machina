package machina

import (
	"fmt"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var stopCommand = &cobra.Command{
	Use:   "stop",
	Short: "Stops a running machine",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		if len(args) > 0 {
			name = args[0]
		}
		machine, err := vm.Load(name)
		if err != nil {
			log.Error("the machine doesn't exist")
		}

		// Start the machine
		log.Info(fmt.Sprintf("Starting machine '%s'", name))
		machine.Stop()
	},
}

func init() {
	rootCommand.AddCommand(stopCommand)
}
