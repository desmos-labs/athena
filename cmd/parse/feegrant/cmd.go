package feegrant

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewFeeGrant returns the Cobra command that allows to parse things related to the x/feegrant module
func NewFeeGrant(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feegrant",
		Short: "Parse things related to the x/feegrant module",
	}

	cmd.AddCommand(
		allowancesCmd(parseCfg),
	)

	return cmd
}
