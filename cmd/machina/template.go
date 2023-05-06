package machina

import (
	"fmt"

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
			tpls := vm.GetTemplateList()

			log.Printf("Available templates")
			for _, tpl := range tpls {
				fmt.Println(fmt.Sprintf("%s", string(tpl)))
			}

			log.Println("Use 'machina template <name>' to get a specific template.")
		}

	},
}

func init() {
	rootCommand.AddCommand(templateCommand)
}
