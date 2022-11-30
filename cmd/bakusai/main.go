package main

import (
	"io"
	"log"

	"github.com/alecthomas/kong"
	"github.com/bonnou-shounen/bakusai/internal/cmd"
)

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(
		&cli,
		kong.Name("bakusai"),
		kong.ShortUsageOnError(),
	)

	if cli.Debug {
		log.SetOutput(ctx.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	ctx.FatalIfErrorf(ctx.Run(&cli))
}
