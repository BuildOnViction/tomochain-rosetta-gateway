package cmd

import (
	"github.com/tomochain/tomochain-rosetta-gateway/tomochain"

	"github.com/spf13/cobra"
)

var (
	utilsBootstrapCmd = &cobra.Command{
		Use:   "utils:generate-bootstrap",
		Short: "Generate a bootstrap balance configuration file",
		Long: `For rosetta-cli testing, it can be useful to generate
a bootstrap balances file for balances that were created
at genesis. This command creates such a file given the
path of an Ethereum genesis file.
When calling this command, you must provide 2 arguments:
[1] the location of the genesis file
[2] the location of where to write bootstrap balances file`,
		RunE: runUtilsBootstrapCmd,
		Args: cobra.ExactArgs(2), //nolint:gomnd
	}
)

func runUtilsBootstrapCmd(cmd *cobra.Command, args []string) error {
	return tomochain.GenerateBootstrapFile(args[0], args[1])
}
