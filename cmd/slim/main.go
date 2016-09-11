package main

import (
	"fmt"
	"os"
  "os/exec"
  "syscall"

  "github.com/ripta/funnel"
)

func main() {
  name := os.Args[0]

  if len(os.Args) < 2 {
    fmt.Fprintf(os.Stderr, "%s: binary name required\n", name)
    // fmt.Fprintf(os.Stderr, "Usage: %s /path/to/binary [binary-arguments...]\n", name)
    os.Exit(-1)
  }

  binary := os.Args[1]
  binaryArgs := os.Args[2:]

  opts := &funnel.FunnelOptions{
    Format:                "json",
    LogTimestamp:          true,
    LogTimestampFieldName: "@timestamp",
    RedirectStderr:        []string{"-"},
    RedirectStdout:        []string{"-"},
  }

  exitCode, err := funnel.Exec(opts, binary, binaryArgs)
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
				os.Exit(int(ws.ExitStatus())) // POSIX
				// os.Exit(int(ws.ExitCode)) // Windows
			}

			fmt.Fprintf(os.Stderr, "%s: unknown status code %q", name, ee)
			os.Exit(-1)
		}

		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(-1)
	}

	os.Exit(exitCode)
}
