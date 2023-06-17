package machina

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexeyco/simpletable"
	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "Lists all created instances",
	Aliases: []string{"ls"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {

		// Load the application configuration
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading configuration")
		}

		type inst struct {
			name    string
			cluster string
		}
		var instances []inst
		// Get's a list of all the instances created
		dirs, _ := os.ReadDir(cfg.Directories.Instances)
		for _, dir := range dirs {
			i := inst{
				name: dir.Name(),
			}
			instances = append(instances, i)
		}

		// Get's a list of all the instances on all the created clusters
		dirs, _ = os.ReadDir(cfg.Directories.Clusters)
		for _, dir := range dirs {
			d := filepath.Join(cfg.Directories.Clusters, dir.Name())
			newDirs, _ := os.ReadDir(d)
			for _, newDir := range newDirs {
				i := inst{
					name:    fmt.Sprintf("%s/%s", dir.Name(), newDir.Name()),
					cluster: dir.Name(),
				}
				instances = append(instances, i)
			}
		}

		// Create a new visual table and set the header titles
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "VM NAME"},
				{Align: simpletable.AlignCenter, Text: "CLUSTER"},
				{Align: simpletable.AlignCenter, Text: "IP ADDRESS"},
				{Align: simpletable.AlignCenter, Text: "STATUS"},
				{Align: simpletable.AlignCenter, Text: "CPUS"},
				{Align: simpletable.AlignCenter, Text: "MEMORY"},
				{Align: simpletable.AlignCenter, Text: "DISK"},
				{Align: simpletable.AlignCenter, Text: "LABELS"},
			},
		}

		// Add the content for all the rows
		for _, instance := range instances {
			kind, _ := hypvsr.Load(instance.name)
			vms := kind.GetVMs()
			for _, vm := range vms {
				status, err := vm.Status()
				if err != nil {
					status = "error"
				}

				r := []*simpletable.Cell{
					{Text: vm.Name},
					{Text: instance.cluster},
					{Text: vm.Network.IPAddress},
					{Text: status},
					{Align: simpletable.AlignCenter, Text: vm.Resources.CPUs},
					{Align: simpletable.AlignCenter, Text: vm.Resources.Memory},
					{Align: simpletable.AlignCenter, Text: vm.Resources.Disk},
					{Text: vm.Variant},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
		}

		// Print the table
		table.SetStyle(simpletable.StyleDefault)
		fmt.Println(table.String())
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
