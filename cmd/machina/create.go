package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/vm"
	"github.com/spf13/cobra"
)

// Creates a new VM
var createCommand = &cobra.Command{
	Use:     "create",
	Short:   "Creates a new machine",
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		name = "default"

		// Check the passed parameters
		switch {
		// Get the template name to load from the first argument,
		// if passed
		case len(args) == 1:
			name = args[0]
		// Get the filename from the file parameter
		case file != "":
			name = file
		}

		// Create a new MachinaVM struct and the necessary files
		fmt.Printf("Creating necessary files\n")
		vm, err := vm.NewVM(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating files\n")
			os.Exit(1)
		}

		// Download the distro image used for the machine
		fmt.Printf("Downloading image\n")
		err = vm.DownloadImage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downoading image\n")
			os.Exit(1)
		}

		// Create boot and seed disks necessary for the machine to boot
		fmt.Printf("Create boot and seed disks\n")
		err = vm.CreateDisks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating boot and seed disks\n")
			os.Exit(1)
		}

		// Create and start the VM
		fmt.Printf("Create and start the VM\n")
		err = vm.Create()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating machine\n")
			os.Exit(1)
		}

		// Wait until the VM reaches running state
		fmt.Printf("Waiting for the machine to become ready\n")
		err = vm.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "The machine appears to be stuck in a starting state\n")
		}

		// run installation scripts
		fmt.Printf("Running install scripts\n")
		err = vm.RunInitScripts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running install scripts\n")
			os.Exit(1)
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the file to use to create the machine")
	rootCommand.AddCommand(createCommand)
}
