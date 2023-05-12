package machina

import (
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/enkodr/machina/internal/config"
	"github.com/enkodr/machina/internal/vm"
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

		cfg := config.LoadConfig()
		dirs, _ := os.ReadDir(cfg.Directories.Instances)

		// Create a new visual table and set the header titles
		table := simpletable.New()
		table.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "NAME"},
				{Align: simpletable.AlignCenter, Text: "IP"},
				{Align: simpletable.AlignCenter, Text: "STATUS"},
				{Align: simpletable.AlignCenter, Text: "CPUS"},
				{Align: simpletable.AlignCenter, Text: "MEMORY"},
				{Align: simpletable.AlignCenter, Text: "DISK"},
				{Align: simpletable.AlignCenter, Text: "VARIANT"},
			},
		}

		// Add the content for all the rows
		for _, dir := range dirs {
			vmc, _ := vm.Load(dir.Name())
			status, err := vmc.Status()
			if err != nil {
				status = "error"
			}
			r := []*simpletable.Cell{
				{Text: vmc.Name},
				{Text: vmc.Network.IPAddress},
				{Text: status},
				{Align: simpletable.AlignCenter, Text: vmc.Specs.CPUs},
				{Align: simpletable.AlignCenter, Text: vmc.Specs.Memory},
				{Align: simpletable.AlignCenter, Text: vmc.Specs.Disk},
				{Text: vmc.Variant},
			}
			table.Body.Cells = append(table.Body.Cells, r)
		}

		// Print the table
		table.SetStyle(simpletable.StyleDefault)
		fmt.Println(table.String())
	},
}

func init() {
	rootCommand.AddCommand(listCommand)
}
