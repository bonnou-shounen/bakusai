package cmd

type optGlobal struct {
	Debug bool `hidden:"" env:"BAKUSAI_DEBUG"`
}

type CLI struct {
	optGlobal
	Dump struct {
		Thread DumpThread `cmd:""`
	} `cmd:""`
	Version PrintVersion `cmd:"" hidden:""`
}
