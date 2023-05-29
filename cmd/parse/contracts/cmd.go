package contracts

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewContractsCmd returns the Cobra command that allows to parse things related to smart contracts
func NewContractsCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contracts",
		Short: "Parse things related to smart contracts",
	}

	cmd.AddCommand(
		contractsCmd(parseCfg),
	)

	return cmd
}
