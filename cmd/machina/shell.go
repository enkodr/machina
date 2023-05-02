package machina

import (
	"fmt"
	"os"

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
		if len(args) > 0 {
			name = args[0]
		}
		machine, err := vm.Load(name)
		if err != nil {
			log.Errorf("the machine %q doesn't exist", name)
			os.Exit(0)
		}
		// Start the machine
		log.Info(fmt.Sprintf("Enter the machine '%s' shell", name))
		machine.Shell()

	},
}

func init() {
	rootCommand.AddCommand(shellCommand)
}
