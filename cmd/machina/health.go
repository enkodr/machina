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
	Short: "Checks if the application is able to run the machines",
	Run: func(cmd *cobra.Command, args []string) {
		// Format the output to show as table

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Printf("Checking if dependencies are installed...\n")

		var deps []string
		if runtime.GOOS == "linux" {
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
			deps = []string{
				"qemu-system-x86_64",
			}
		}
		for _, dep := range deps {
			if packageInstalled(dep) {
				fmt.Fprintln(w, fmt.Sprintf("%s\tinstalled", dep))
			} else {
				fmt.Fprintln(w, fmt.Sprintf("%s\tnot installed", dep))
			}
		}
		w.Flush()
	},
}

func init() {
	rootCommand.AddCommand(healthCommand)
}

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
