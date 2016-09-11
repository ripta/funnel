package main

import (
	// "fmt"
	// "io/ioutil"
	// "os"

	"github.com/ripta/funnel"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func funnelOptionsFromContext(c *cobra.Command) *funnel.FunnelOptions {
	opts := &funnel.FunnelOptions{
		Format:                viper.GetString("format"),
		LogTimestamp:          viper.GetBool("log-timestamp"),
		LogTimestampFieldName: viper.GetString("log-timestamp-field-name"),
		RedirectStderr:        shimBrokenStringSlice(c, "redirect-stderr"),
		RedirectStdout:        shimBrokenStringSlice(c, "redirect-stdout"),
	}

	return opts
}

// viper.BindPFlag(s) borks pflag.StringSlice{} values by stringifying the
// pflag value, and then wrapping it in a string slice. This shim tries
// pflag.GetStringSlice first. If the call fails, or if the value is empty,
// only then do we ask viper.
//
// See also https://github.com/spf13/viper/issues/112 for reference.
func shimBrokenStringSlice(c *cobra.Command, name string) []string {
	values, err := c.PersistentFlags().GetStringSlice(name)
	if err == nil && values != nil && len(values) > 0 {
		return values
	}

	return viper.GetStringSlice(name)
}
