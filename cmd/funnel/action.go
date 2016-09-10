package main

import (
	"fmt"

	"github.com/ripta/funnel"
	"github.com/urfave/cli"
)

func funnelCliAction(c *cli.Context) error {
	if c.NArg() == 0 {
		msg := fmt.Sprintf("%s: at least one argument (the name of the binary) is required", c.App.Name)
		return cli.NewExitError(msg, 255)
	}

  fmt.Printf("redirect-stderr = %v\n", c.StringSlice("redirect-stderr"))
  fmt.Printf("redirect-stdout = %v\n", c.StringSlice("redirect-stdout"))
	return funnel.Exec(funnelOptionsFromContext(c), c.Args().First(), c.Args().Tail())
}

func funnelOptionsFromContext(c *cli.Context) *funnel.FunnelOptions {
	opts := &funnel.FunnelOptions{}
	return opts
}
