package tomochain

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/sync/errgroup"
)

const (
	tomoLogger       = "tomo"
	tomoStdErrLogger = "tomochain log"
)

// logPipe prints out logs from geth. We don't end when context
// is canceled beacause there are often logs printed after this.
func logPipe(pipe io.ReadCloser, identifier string) error {
	reader := bufio.NewReader(pipe)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println("closing", identifier, err)
			return err
		}

		message := strings.ReplaceAll(str, "\n", "")
		log.Println(identifier, message)
	}
}

// StartTomo starts a geth daemon in another goroutine
// and logs the results to the console.
func StartTomo(ctx context.Context, arguments string, g *errgroup.Group) error {
	parsedArgs := strings.Split(arguments, " ")

	// get datadir
	// default datadir = /data
	datadir := "/data"
	for _, arg := range parsedArgs {
		if strings.HasPrefix(arg, "--datadir") {
			d := strings.Split(arg, "=")
			if len(d) > 0 {
				datadir = d[1]
			}
		}
	}
	if _, err := os.Stat(datadir); os.IsNotExist(err) {
		fmt.Println("create data dir", datadir)
		os.Mkdir(datadir, 755)
	}

	// initialize if not exist
	if _, err := os.Stat(path.Join(datadir, "tomo")); os.IsNotExist(err) {
		fmt.Println("Initialize tomochain datadir with genesis")
		initCmd := exec.Command(
			"/app/tomo",
			"init",
			"/app/genesis.json",
			"--datadir="+datadir,
		)
		if err := initCmd.Run(); err != nil {
			fmt.Println("Failed to initialize tomochain datadir", err)
		}
	}

	cmd := exec.Command(
		"/app/tomo",
		parsedArgs...,
	) // #nosec G204

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	g.Go(func() error {
		return logPipe(stdout, tomoLogger)
	})

	g.Go(func() error {
		return logPipe(stderr, tomoStdErrLogger)
	})

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%w: unable to start tomo", err)
	}

	g.Go(func() error {
		<-ctx.Done()

		log.Println("sending interrupt to tomo")
		return cmd.Process.Signal(os.Interrupt)
	})

	return cmd.Wait()
}
