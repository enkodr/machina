package main

import "github.com/enkodr/machina/cmd/machina"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	machina.Execute()
}
