package main

import (
	"fmt"
	"os"
	// "github.com/spf13/cobra"
)

func main() {
	// cobra.OnInitialize(init)
	if err := funnelRoot().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}
