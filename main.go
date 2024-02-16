package main

import "github.com/enkodr/machina/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.Execute()
}
