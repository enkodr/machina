package machina

import (
	"fmt"
	"os"
	"strings"

	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var copyCommand = &cobra.Command{
	Use:   "copy",
	Short: "Copies files/directories from the host to the VM and vice-versa",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			log.Error("Error!\n")
			log.Info("\tmachina copy <host_path> <machine_name>:<machine_path>\n")
			log.Info("\tmachina copy <machine_name>:<machine_path> <host_path>\n")
			os.Exit(0)
		}
		return nil
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		if len(args) > 0 {
			name = args[0]
		}

		// Identify the origin from copy (Host or VM)
		hostToVM := true
		origin := args[0]
		dest := args[1]
		if strings.Contains(origin, ":") {
			name = strings.Split(origin, ":")[0]
			hostToVM = false
		} else {
			name = strings.Split(dest, ":")[0]
		}

		// Get the machine
		machine, err := vm.Load(name)
		if err != nil {
			log.Errorf("the machine %q doesn't exist", name)
			os.Exit(0)
		}

		// Defiine the origin and destination
		if hostToVM {
			parts := strings.Split(dest, ":")
			dest = fmt.Sprintf("%s@%s:%s", machine.Credentials.Username, machine.Network.IPAddress, parts[1])
		} else {
			parts := strings.Split(origin, ":")
			origin = fmt.Sprintf("%s@%s:%s", machine.Credentials.Username, machine.Network.IPAddress, parts[1])
		}

		// Copy the content
		log.Info("Copying content... ")

		// Start the machine
		machine.CopyContent(origin, dest)

		// Wait until machine is ready
		machine.Wait()

		log.Info("Done!\n")

	},
}

func init() {
	rootCommand.AddCommand(copyCommand)
}
