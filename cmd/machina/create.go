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
	Short:   "Creates a new instance",
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {
		name = "default"

		// Evaluate the passed parameters to determine the template name
		switch {
		// If exactly one argument has been provided, it is assumed to be the name of the template.
		// Hence, assign this first argument to the variable 'name'.
		case len(args) == 1:
			name = args[0]
		// If the 'file' parameter is not an empty string, it is considered to be the filename from which the template should be loaded.
		// Therefore, assign the 'file' parameter value to the variable 'name'.
		case file != "":
			name = file
		}

		fmt.Printf("Creating instance\n")
		// Call the NewInstance function of the hypvsr package to create a new instance instance with the given name.
		// The function returns a reference to the new instance instance and any error that may have occurred.
		instance, err := hypvsr.NewInstance(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance\n")
			os.Exit(1)
		}

		fmt.Printf("Creating necessary files\n")
		// Call the CreateDir method that will create the directory where the instance will be created
		err = instance.CreateDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Instance already exists\n")
			os.Exit(1)
		}
		// Call the Prepare method that will create all the necessary files needed for the instance to work
		err = instance.Prepare()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating files\n")
			os.Exit(1)
		}

		fmt.Printf("Downloading image\n")
		// Call the DownloadImage method that will download the distro image needed to boot the instance
		err = instance.DownloadImage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downoading image\n")
			os.Exit(1)
		}

		// Call the CreateDisks methiod that will create boot and seed disks necessary for the instance to boot
		fmt.Printf("Create boot and seed disks\n")
		err = instance.CreateDisks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating boot and seed disks\n")
			os.Exit(1)
		}

		fmt.Printf("Create and start the VM\n")
		// Call the Create method that will create and start the instance
		err = instance.Create()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance\n")
			os.Exit(1)
		}

		fmt.Printf("Waiting for the instance to become ready\n")
		// Call the Wait methid that will wait until the VM reaches running state
		err = instance.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "The instance appears to be stuck in a starting state\n")
		}

		fmt.Printf("Running install scripts\n")
		// Call the RunInitScripts method that will run the initial scripts defined on the template
		err = instance.RunInitScripts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running install scripts\n")
			os.Exit(1)
		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the file to use to create the instance")
	createCommand.PersistentFlags().StringVarP(&newName, "name", "n", "", "specify the name of the instance")
	createCommand.PersistentFlags().StringVarP(&newCPUs, "cpus", "c", "", "specify the amount of CPUs of the instance (e.g. 2)")
	createCommand.PersistentFlags().StringVarP(&newMem, "mem", "m", "", "specify the amount of memory of the instance (e.g 4G)")
	createCommand.PersistentFlags().StringVarP(&newDisk, "disk", "d", "", "specify the size of the disk of the instance (e.g. 100G)")
	rootCommand.AddCommand(createCommand)
}
