package main

import (
	// "fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func tryLoadingYamlFromFlag(flag, defaultFile string) cli.BeforeFunc {
	return func(c *cli.Context) error {
		var filename string

		if c.IsSet(flag) {
			if _, err := os.Stat(c.String(flag)); os.IsNotExist(err) {
				return err
			}

			filename = c.String(flag)
		}

		if defaultFile != "" {
			if _, err := os.Stat(defaultFile); !os.IsNotExist(err) {
				filename = defaultFile
			}
		}

		if filename != "" {
			var results map[string]string
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				return err
			}

			yaml.Unmarshal(content, &results)

			for k, v := range results {
				if c.IsSet(k) {
					continue
				}

				c.Set(k, v)

				// FIXME(rpasay): Looks like c.IsSet is actually cached, so we can't easily...
				// if !c.IsSet(k) {
				// 	return cli.NewExitError(fmt.Sprintf("Invalid configuration key '%v' was set in '%s'", k, filename), 253)
				// }
			}
		}

		return nil
	}
}
