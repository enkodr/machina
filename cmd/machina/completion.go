package machina

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/config"
	"github.com/spf13/cobra"
)

// Get's the commands and the created machines names for auto-completion
func bashCompleteInstanceNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Load configuration from file
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Get instances names for completion
	var instances []string
	vms, err := os.ReadDir(cfg.Directories.Machines)
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
	// Get instances names
	instances, _ := bashCompleteInstanceNames(cmd, args, toComplete)

	// Add :/ to each instance completion to make it easier for
	// defining paths inside the machine
	for i := 0; i < len(instances); i++ {
		instances[i] = fmt.Sprintf("%s:/", instances[i])
	}
	return instances, cobra.ShellCompDirectiveNoSpace
}
