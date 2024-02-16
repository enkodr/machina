package cmd

import (
	"fmt"
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/enkodr/machina/internal/hypvsr"
	"github.com/spf13/cobra"
)

var templateCommand = &cobra.Command{
	Use:     "template",
	Short:   "Lists and gets the available templates",
	Aliases: []string{"tpl"},
	Args: func(cmd *cobra.Command, args []string) error {
		return validateName(cmd, args)
	},
	ValidArgsFunction: bashCompleteInstanceNames,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the instance data

		// Check if a template name was passed
		if len(args) == 1 {
			tpl, err := hypvsr.GetTemplate(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to get template %s\n", args[0])
				os.Exit(1)
			}
			fmt.Println(tpl)
		} else {
			fmt.Println("Available templates")
			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Align: simpletable.AlignCenter, Text: "TEMPLATE NAME"},
				},
			}
			// Get the template list
			tpls := hypvsr.GetTemplateList()

			for _, tpl := range tpls {
				r := []*simpletable.Cell{
					{Text: tpl},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}

			table.SetStyle(simpletable.StyleDefault)
			fmt.Println(table.String())

			fmt.Printf("Use 'machina template <name>' to get a specific template.\n")
		}

	},
}

func init() {
	rootCommand.AddCommand(templateCommand)
}
