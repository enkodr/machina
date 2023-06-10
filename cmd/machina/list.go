package machina

import (
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "Lists all created machines",
	Aliases: []string{"ls"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading configuration")
		}
		dirs, _ := os.ReadDir(cfg.Directories.Machines)

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
		for _, dir := range dirs {
			kind, _ := hypvsr.Load(dir.Name())
			vms := kind.GetVMs()
			for _, vm := range vms {
				status, err := vm.Status()
				if err != nil {
					status = "error"
				}

				r := []*simpletable.Cell{
					{Text: vm.Name},
					{Text: ""},
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
