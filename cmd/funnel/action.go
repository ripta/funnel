package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/ripta/funnel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	jww "github.com/spf13/jwalterweatherman"
)

func funnelRoot() *cobra.Command {
	if os.Getenv("FUNNEL_TRACE") != "" {
		jww.SetStdoutThreshold(jww.LevelTrace)
	}

	root := &cobra.Command{
		Use:   "funnel",
		Short: "Execute a binary and redirect its output somewhere",
		Run:   funnelCliAction,
	}

	viper.SetConfigName("funnel")
	viper.AddConfigPath("/etc/funnel/")
	viper.AddConfigPath("$HOME/.funnel")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("funnel")
	viper.AutomaticEnv()

	pf := root.PersistentFlags()
	pf.StringP("config", "c", "", "Configuration file in .json or .yml (default is $HOME/.funnel.json, fallback /etc/funnel/funnel.json)")
	pf.StringP("format", "f", "json", "Log format: json, raw (default is json)")
	pf.Bool("log-timestamp", true, "Whether to include timestamp (as a field in JSON or prefix in raw)")
	pf.String("log-timestamp-field-name", "@timestamp", "(JSON only) The name of the timestamp field")
	pf.StringSlice("redirect-stderr", nil, "One or more locations to redirect STDERR into (specify multiple times)")
	pf.StringSlice("redirect-stdout", nil, "One or more locations to redirect STDOUT into (specify multiple times)")

	viper.BindPFlags(pf)
	viper.ReadInConfig()

	return root
}

func funnelCliAction(c *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "%s: at least one argument (the name of the binary) is required\n", c.Name())
		os.Exit(-1)
	}

	exitCode, err := funnel.Exec(funnelOptionsFromContext(c), args[0], args[1:])
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
				os.Exit(int(ws.ExitStatus())) // POSIX
				// os.Exit(int(ws.ExitCode)) // Windows
			}

			fmt.Fprintf(os.Stderr, "%s: unknown status code %q", c.Name(), ee)
			os.Exit(-1)
		}

		fmt.Fprintf(os.Stderr, "%s: %v\n", c.Name(), err)
		os.Exit(-1)
	}

	os.Exit(exitCode)
}
