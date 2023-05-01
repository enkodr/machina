package machina

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists all created machines",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.LoadConfig()
		dirs, _ := os.ReadDir(cfg.Directories.Instances)

		log.Info("List all machines")
		w := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\t IP\t Status\t CPUs\t Memory\t Disk\t Variant")
		for _, dir := range dirs {
			vmc, _ := vm.Load(dir.Name())
			status, err := vmc.Status()
			if err != nil {
				status = "error"
			}
			fmt.Fprintln(w,
				vmc.Name,
				"\t", vmc.Network.IPAddress,
				"\t", status,
				"\t", vmc.Specs.CPUs,
				"\t", vmc.Specs.Memory,
				"\t", vmc.Specs.Disk,
				"\t", vmc.Variant,
			)
		}
		w.Flush()
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
