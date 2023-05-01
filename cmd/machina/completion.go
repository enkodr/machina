package machina

import (
	"io/ioutil"

	"github.com/enkodr/machina/internal/config"
	"github.com/spf13/cobra"
)

// Get's the commands and the created machines names for auto-completion
func bashCompleteInstanceNames(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	cfg := config.LoadConfig()

	// Get instances names for completion
	instances := []string{}
	vms, err := ioutil.ReadDir(cfg.Directories.Instances)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	for _, vm := range vms {
		instances = append(instances, vm.Name())
	}
	return instances, cobra.ShellCompDirectiveNoFileComp
}
