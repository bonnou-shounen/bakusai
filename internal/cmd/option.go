package cmd

type Option struct {
	Debug      bool         `hidden:"" env:"BAKUSAI_DEBUG"`
	DumpThread DumpThread   `cmd:""`
	Version    PrintVersion `cmd:"" hidden:""`
}
