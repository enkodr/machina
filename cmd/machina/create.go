package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var (
	newName string
	newCPUs string
	newMem  string
	newDisk string
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
		// Get the template name to load from the first argument, if passed
		case len(args) == 1:
			name = args[0]
		// Get the filename from the file parameter
		case file != "":
			name = file
		}

		// Create a new VM
		fmt.Printf("Creating machine\n")
		vm, err := hypvsr.NewInstance(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating machine\n")
			os.Exit(1)
		}

		// Prepare necessary files for machine creation
		fmt.Printf("Creating necessary files\n")
		err = vm.CreateDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Instance already exists\n")
			os.Exit(1)
		}
		err = vm.Prepare()
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
	createCommand.PersistentFlags().StringVarP(&newName, "name", "n", "", "specify the name of the machine")
	createCommand.PersistentFlags().StringVarP(&newCPUs, "cpus", "c", "", "specify the amount of CPUs of the machine (e.g. 2)")
	createCommand.PersistentFlags().StringVarP(&newMem, "mem", "m", "", "specify the amount of memory of the machine (e.g 4G)")
	createCommand.PersistentFlags().StringVarP(&newDisk, "disk", "d", "", "specify the size of the disk of the machine (e.g. 100G)")
	rootCommand.AddCommand(createCommand)
}
