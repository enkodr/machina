package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/config"
	"github.com/spf13/cobra"
)

// Get's the commands and the created machines names for auto-completion
func bashCompleteInstanceNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg := config.LoadConfig()

	// Get instances names for completion
	instances := []string{}
	vms, err := os.ReadDir(cfg.Directories.Instances)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	for _, vm := range vms {
		instances = append(instances, vm.Name())
	}
	return instances, cobra.ShellCompDirectiveNoFileComp
}

// Get's the existing VM list if a format easier to complete
func bashCompleteInstanceNamesConnection(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	instances, _ := bashCompleteInstanceNames(cmd, args, toComplete)
	size := len(instances)
	for i := 0; i < size; i++ {
		instances[i] = fmt.Sprintf("%s:/", instances[i])
	}
	return instances, cobra.ShellCompDirectiveNoSpace
}
