package machina

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		dirs, err := os.ReadDir(cfg.Directories.Instances)
		if err != nil {
			log.Error("failed to get machines")
		}

		log.Info("List all machines")
		w := tabwriter.NewWriter(os.Stdout, 10, 0, 2, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Name\tIP\tStatus\tCPUs\tMemory\tDisk\tVariant")
		for _, dir := range dirs {
			data, _ := os.ReadFile(filepath.Join(cfg.Directories.Instances, dir.Name(), "machina.yaml"))
			vm := &vm.VMConfig{}
			yaml.Unmarshal(data, vm)

			fmt.Fprintln(w,
				vm.Name,
				"\t", vm.Network.IPAddress,
				"\t", "",
				"\t", vm.Specs.CPUs,
				"\t", vm.Specs.Memory,
				"\t", vm.Specs.Disk,
				"\t", vm.Variant,
			)
		}
		w.Flush()
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
