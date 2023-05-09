package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
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
		if len(args) > 0 {
			name = args[0]
		}
		machine, err := vm.Load(name)
		if err != nil {
			log.Errorf("the machine %q doesn't exist", name)
			os.Exit(0)
		}

		// Start the machine
		log.Info(fmt.Sprintf("Starting machine '%s'", name))
		machine.Start()

		// Wait until machine is ready
		log.Info(fmt.Sprintf("Waiting for machine '%s'", name))
		machine.Wait()

	},
}

func init() {
	rootCommand.AddCommand(startCommand)
}
