package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var shellCommand = &cobra.Command{
	Use:     "shell",
	Short:   "Enters a machine shell",
	Aliases: []string{"sh"},
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
		// Entering the machine shell
		log.Info(fmt.Sprintf("Enter machine's '%s' shell", name))
		err = machine.Shell()
		if err != nil {
			log.Errorf("error entering machine's %q shell", name)
			os.Exit(1)
		}
	},
}

func init() {
	rootCommand.AddCommand(shellCommand)
}
