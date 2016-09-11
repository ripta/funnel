package funnel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type LogTransformerFunc func(src []byte, message []byte) ([]byte, error)

type FunnelExecutor struct{}
type FunnelOptions struct {
	Format                string
	LogTimestamp          bool
	LogTimestampFieldName string
	RedirectStderr        []string
	RedirectStdout        []string
}

func (o *FunnelOptions) BuildTransformer() (LogTransformerFunc, error) {
	switch o.Format {
	case "json":
		return jsonLogTransformer, nil
	case "raw":
		return rawLogTransformer, nil
	default:
		return nil, fmt.Errorf("unrecognized format option %q", o.Format)
	}
}

func Exec(opts *FunnelOptions, command string, commandArgs []string) (int, error) {
	logTransformer, err := opts.BuildTransformer()
	if err != nil {
		return -1, err
	}

	cmd := exec.Command(command, commandArgs...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}

	// fileMode := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	// errFile, err := os.OpenFile("stderr.json", fileMode, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer errFile.Close()

	// outFile, err := os.OpenFile("stdout.json", fileMode, 0644)
	// if err != nil {
	// 	return err
	// }
	// defer outFile.Close()

	errFile := os.Stdout
	outFile := os.Stdout

	workers := sync.WaitGroup{}
	dispatchBufferWorker(&workers, logTransformer, []byte("stderr"), errFile, stderr)
	dispatchBufferWorker(&workers, logTransformer, []byte("stdout"), outFile, stdout)

	if err = cmd.Start(); err != nil {
		return -1, err
	}

	// Workers must finish reading from the pipes before we can cmd.Wait()
	workers.Wait()

	// The caller can optionally handle specific error codes themselves;
	// see `cmd/funnel/action.go` for an example
	if err = cmd.Wait(); err != nil {
		return -1, err
	}

	return 0, nil
}

func dispatchBufferWorker(workers *sync.WaitGroup, logTransformer LogTransformerFunc, src []byte, w io.Writer, r io.ReadCloser) {
	s := bufio.NewScanner(r)
	workers.Add(1)

	go func() {
		defer workers.Done()

		for s.Scan() {
			msg, err := logTransformer(src, s.Bytes())
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			_, err = w.Write(msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				break
			}

			_, err = w.Write([]byte{0x0A}) // 0x0A == "\n"
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				break
			}
		}
	}()
}
