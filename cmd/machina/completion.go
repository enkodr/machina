package machina

import (
	"github.com/spf13/cobra"
)

// Get's the commands and the created machines names for auto-completion
func bashCompleteInstanceNames(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	// cfg, err := config.New()
	// // Get shell completion from cobra
	// if err != nil {
	// 	return nil, cobra.ShellCompDirectiveDefault
	// }
	// // Get instances names for completion
	// instances := []string{}
	// vms, err := ioutil.ReadDir(cfg.InstancesDirectory)
	// if err != nil {
	// 	return nil, cobra.ShellCompDirectiveDefault
	// }
	// for _, vm := range vms {
	// 	instances = append(instances, vm.Name())
	// }
	// return instances, cobra.ShellCompDirectiveNoFileComp
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}
