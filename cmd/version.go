package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	gitCommit = "none"
	buildDate = "unknown"
)

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Display application version",
	Run: func(cmd *cobra.Command, args []string) {
		// Outputs the parsed yaml string into the screen
		fmt.Printf("version: %s\n", version)
		fmt.Printf("git commit: %s\n", gitCommit)
		fmt.Printf("build date: %s\n", buildDate)
	},
}

func init() {
	rootCommand.AddCommand(versionCommand)
}
