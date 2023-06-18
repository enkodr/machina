package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
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

		fmt.Printf("Creating instance %s\n", name)
		// Call the NewInstance function of the hypvsr package to create a new instance instance with the given name.
		// The function returns a reference to the new instance instance and any error that may have occurred.
		tpl := hypvsr.NewTemplate(name)
		instance, err := hypvsr.NewInstance(tpl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating instance\n")
			os.Exit(1)
		}

		fmt.Printf("Creating necessary files\n")
		// Call the CreateDir method that will create the directory where the instance will be created
		for _, machine := range instance.Machines {
			err = machine.CreateDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Instance already exists\n")
				os.Exit(1)
			}

			// Call the Prepare method that will create all the necessary files needed for the instance to work
			err = machine.Prepare()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating files\n")
				os.Exit(1)
			}

			fmt.Printf("Downloading image\n")
			// Call the DownloadImage method that will download the distro image needed to boot the instance
			err = machine.DownloadImage()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading images\n")
				os.Exit(1)
			}

			fmt.Printf("Installing instance %q\n", machine.Name)
			// Call the CreateDisks methiod that will create boot and seed disks necessary for the instance to boot
			fmt.Printf("Create boot and seed disks for instance '%s'\n", machine.Name)
			err = machine.CreateDisks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating boot and seed disks\n")
				os.Exit(1)
			}

			fmt.Printf("Create and start instance %q\n", machine.Name)
			// Call the Create method that will create and start the instance
			err = machine.Create()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating instance\n")
				os.Exit(1)
			}

			fmt.Printf("Waiting for the instance %q to become ready\n", machine.Name)
			// Call the Wait methid that will wait until the VM reaches running state
			err = machine.Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "The instance appears to be stuck in a starting state\n")
			}

			fmt.Printf("Running install scripts in instance %q\n", machine.Name)
			// Call the RunInitScripts method that will run the initial scripts defined on the template
			err = machine.RunInitScripts()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running install scripts in instance\n")
				os.Exit(1)
			}

		}

		fmt.Printf("Done!\n")
	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the file to use to create the instance")
	rootCommand.AddCommand(createCommand)
}
