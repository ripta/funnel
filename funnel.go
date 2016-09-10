package funnel

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

type FunnelExecutor struct{}
type FunnelOptions struct {
	ConfigFilename string
}

type LogTransformerFunc func(src []byte, message []byte) ([]byte, error)

func Exec(opts *FunnelOptions, command string, commandArgs []string) error {
	cmd := exec.Command(command, commandArgs...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
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
	logTransformer := jsonLogTransformer

	workers := sync.WaitGroup{}
	dispatchBufferWorker(&workers, logTransformer, []byte("stderr"), errFile, stderr)
	dispatchBufferWorker(&workers, logTransformer, []byte("stdout"), outFile, stdout)

	if err = cmd.Start(); err != nil {
		return err
	}

	// Workers must finish reading from the pipes before we can cmd.Wait()
	workers.Wait()

	// fmt.Fprintf(os.Stderr, "main done\n")
	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
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
