package machina

import (
	"fmt"
	"io/fs"
	"os"

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

		var instances []fs.DirEntry
		// Get's a list of all the instances created
		dirs, _ := os.ReadDir(cfg.Directories.Instances)
		instances = append(instances, dirs...)

		// Create a new visual table and set the header titles
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "VM NAME"},
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
			kind, _ := hypvsr.GetMachine(instance.Name())
			vms := kind.GetVMs()
			for _, vm := range vms {
				status, err := vm.Status()
				if err != nil {
					status = "error"
				}

				r := []*simpletable.Cell{
					{Text: vm.Name},
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
