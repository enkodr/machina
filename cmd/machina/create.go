package machina

import (
	"os"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Creates a new VM
var createCommand = &cobra.Command{
	Use:     "create",
	Short:   "Creates a new machine",
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		name = "default"
		if len(args) == 1 {
			name = args[0]
		}

		if file != "" {
			name = file
		}
		// create a new MachinaVM instance
		log.Info("Creating necessary files")
		vm, err := vm.NewVM(name)
		if err != nil {
			log.Error("error creating files")
			os.Exit(1)
		}

		// download the distro image used for the machine
		log.Info("Downloading image")
		err = vm.DownloadImage()
		if err != nil {
			log.Error("error downoading image")
			os.Exit(1)
		}

		// create boot and seed disks
		log.Info("Create boot and seed disks")
		err = vm.CreateDisks()
		if err != nil {
			log.Error("error creating boot and seed disks")
			os.Exit(1)
		}

		// create the VM
		log.Info("Create and start the VM")
		err = vm.Create()
		if err != nil {
			log.Error("error creating machine")
			os.Exit(1)
		}

		// wait until the VM reaches running state
		log.Info("Waiting for the machine to become ready")
		vm.Wait()

		// run installation scripts
		log.Info("Running install scripts")
		err = vm.RunInitScripts()
		if err != nil {
			log.Error("error running install scripts")
			os.Exit(1)
		}

	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the file to use to create the machine")
	rootCommand.AddCommand(createCommand)
}
