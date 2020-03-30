package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.ibm.com/guy-barden/topic-migration/pkg/cmd"
)

type ()

var cli struct {
	Migrate cmd.Migrate `cmd:"" help:"Migrate topic definitions from a source cluster to a target cluster"`
}

func main() {
	sourceAPIKey := os.Getenv("SOURCE_API_KEY")
	targetAPIKey := os.Getenv("TARGET_API_KEY")
	if sourceAPIKey == "" || targetAPIKey == "" {
		fmt.Println("The SOURCE_API_KEY and TARGET_API_KEY environment variables must be set")
		os.Exit(1)
	}

	ctx := kong.Parse(&cli)
	err := ctx.Run(&cmd.Globals{
		SourceAPIKey: sourceAPIKey,
		TargetAPIKey: targetAPIKey,
	})
	ctx.FatalIfErrorf(err)
}
