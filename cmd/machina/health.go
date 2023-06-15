package machina

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var healthCommand = &cobra.Command{
	Use:   "health",
	Short: "Checks if the application is able to run the instances",
	Run: func(cmd *cobra.Command, args []string) {
		// Format the output to show as table

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Printf("Checking if dependencies are installed...\n")

		var deps []string

		// Identify the OS and define the dependencies per OS
		if runtime.GOOS == "linux" {
			// Dependencies for Linux
			deps = []string{
				"cloud-localds",
				"genisoimage",
				"qemu-img",
				"qemu-system-x86_64",
				"scp",
				"ssh",
				"virt-install",
				"virsh",
			}
		} else {
			// Dependenciesfor MacOS
			deps = []string{
				"qemu-system-x86_64",
			}
		}

		// Check if the dependencies are installed and show the outcome
		for _, dep := range deps {
			if packageInstalled(dep) {
				fmt.Fprintln(w, fmt.Sprintf("%s\tinstalled", dep))
			} else {
				fmt.Fprintln(w, fmt.Sprintf("%s\tnot installed", dep))
			}
		}
		err := w.Flush()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error presenting information")
			os.Exit(1)
		}
	},
}

func init() {
	rootCommand.AddCommand(healthCommand)
}

// Check if a dpendency is installed
func packageInstalled(pkg string) bool {
	_, err := exec.LookPath(pkg)

	// check error
	if err != nil {
		// the executable is not found, return false
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return false
		}
	}

	return true
}
