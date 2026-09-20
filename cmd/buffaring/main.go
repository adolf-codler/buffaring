package main

import (
	"log"

	"github.com/alecthomas/kong"

	"adolf-codler/buffaring/internal/cli"
)

func main() {
	var app cli.CLI
	ctx := kong.Parse(&app,
		kong.Name("buffaring"),
		kong.Description("Transfer a buffer over TCP via UDP discovery."),
	)
	if err := ctx.Run(); err != nil {
		log.Fatal(err)
	}
}
