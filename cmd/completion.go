package cmd

import (
	"fmt"
	"os"

	"github.com/enkodr/machina/internal/config"
	"github.com/spf13/cobra"
)

// Get's the commands and the created instances names for auto-completion
func bashCompleteInstanceNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Load configuration from file
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	// Get instanceList names for completion
	var instanceList []string
	instances, err := os.ReadDir(cfg.Directories.Instances)
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	for _, instance := range instances {
		instanceList = append(instanceList, instance.Name())
	}

	// Return the list of instances created
	return instanceList, cobra.ShellCompDirectiveNoFileComp
}

// Get's the existing VM list if a format easier to complete
func bashCompleteInstanceNamesConnection(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Get instances names
	instances, _ := bashCompleteInstanceNames(cmd, args, toComplete)

	// Add :/ to each instance completion to make it easier for
	// defining paths inside the instance
	for i := 0; i < len(instances); i++ {
		instances[i] = fmt.Sprintf("%s:/", instances[i])
	}
	return instances, cobra.ShellCompDirectiveNoSpace
}
