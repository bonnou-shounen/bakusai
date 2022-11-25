package main

import (
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

	ctx.FatalIfErrorf(ctx.Run(&cli))
}
