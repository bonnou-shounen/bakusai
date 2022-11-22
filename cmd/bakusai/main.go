package main

import (
	"github.com/alecthomas/kong"
	"github.com/bonnou-shounen/bakusai/internal/cmd"
)

func main() {
	opt := cmd.Option{}
	ctx := kong.Parse(
		&opt,
		kong.Name("bakusai"),
		kong.ShortUsageOnError(),
	)

	ctx.FatalIfErrorf(ctx.Run(&opt))
}
