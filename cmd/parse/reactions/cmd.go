package profiles

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewReactionsCmd returns the Cobra command that allows to parse things related to the x/reactions module
func NewReactionsCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reactions",
		Short: "Parse things related to the x/reactions module",
	}

	cmd.AddCommand(
		registeredReactionsCmd(parseCfg),
		reactionsCmd(parseCfg),
		paramsCmd(parseCfg),
	)

	return cmd
}
