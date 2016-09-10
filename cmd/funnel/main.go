package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "Execute a binary and redirect its output somewhere"
	app.Version = "0.1.0"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "c, config",
			Value:  "",
			Usage:  "Configuration file from which to load defaults",
			EnvVar: "FUNNEL_CONFIG",
		},
    cli.IntFlag{
      Name:  "verbose",
      Value: 0,
      Usage: "Verbosity level (0-6)",
      EnvVar: "FUNNEL_VERBOSE",
    },
	}

	app.Before = tryLoadingYamlFromFlag("config", "/etc/funnel.yml")
	app.Action = funnelCliAction
	app.Run(os.Args)
}
