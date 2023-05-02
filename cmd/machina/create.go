package machina

import (
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
			log.Fatal(err.Error())
		}

		// download the image
		log.Info("Downloading image")
		err = vm.DownloadImage()
		if err != nil {
			log.Fatal(err.Error())
		}

		// create disks
		log.Info("Create boot disk")
		err = vm.CreateDisks()
		if err != nil {
			log.Fatal(err.Error())
		}

		// create VM
		log.Info("Create and start the VM")
		err = vm.Create()
		if err != nil {
			log.Fatal(err.Error())
		}

		// wait until is running
		log.Info("Waiting for the machine to become ready")
		vm.Wait()

		// run scritps
		log.Info("Running install scripts")
		err = vm.RunInitScripts()
		if err != nil {
			log.Fatal(err.Error())
		}

	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the file to use to create the machine")
	rootCommand.AddCommand(createCommand)
}
