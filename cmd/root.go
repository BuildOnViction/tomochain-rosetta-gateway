package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tomochain-rosetta",
		Short: "TomoChain implementation of the Rosetta API",
	}

	// SignalReceived is set to true when a signal causes us to exit. This makes
	// determining the error message to show on exit much more easy.
	SignalReceived = false
)

// Execute handles all invocations of the
// tomochain-rosetta cmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(utilsBootstrapCmd)
}

// handleSignals handles OS signals so we can ensure we close database
// correctly. We call multiple sigListeners because we
// may need to cancel more than 1 context.
func handleSignals(listeners []context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		color.Red("received signal", "signal", sig)
		SignalReceived = true
		for _, listener := range listeners {
			listener()
		}
	}()
}
