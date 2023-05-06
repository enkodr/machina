package machina

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/enkodr/machina/internal/vm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var templateCommand = &cobra.Command{
	Use:   "template",
	Short: "Lists and gets the available templates",
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: namesBashCompletion,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the machine data
		var tplName string
		if len(args) > 0 {
			tplName = args[0]
		}

		if tplName != "" {
			tpl, err := vm.GetTemplate(tplName)
			if err != nil {
				log.Errorf("Failed to get template %s", tplName)
			}
			fmt.Println(tpl)
		} else {
			log.Info("Available templates")
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "TEMPLATE NAME"},
				},
			}
			tpls := vm.GetTemplateList()

			for _, tpl := range tpls {
				r := []*simpletable.Cell{
					{Text: tpl},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleDefault)
			fmt.Println(table.String())

			log.Println("Use 'machina template <name>' to get a specific template.")
		}

	},
}

func init() {
	rootCommand.AddCommand(templateCommand)
}
