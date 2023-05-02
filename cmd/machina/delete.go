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
		if len(args) > 0 {
			name = args[0]
		}
		vm, err := vm.Load(name)
		if err != nil {
			log.Errorf("the machine %q doesn't exist", name)
			return
		}
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
		// Start the machine
		log.Info(fmt.Sprintf("Deleting machine '%s'", name))
		vm.Delete()

	},
}

func init() {
	rootCommand.AddCommand(deleteCommand)
}
