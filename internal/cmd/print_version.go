package cmd

import (
	"fmt"
	"os"

	"github.com/bonnou-shounen/bakusai"
)

type PrintVersion struct{}

func (*PrintVersion) Run() error {
	fmt.Fprintf(os.Stdout, "%s\n", bakusai.Version)

	return nil
}
