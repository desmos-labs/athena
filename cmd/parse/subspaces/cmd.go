package profiles

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewSubspacesCmd returns the Cobra command that allows to parse things related to the x/subspaces module
func NewSubspacesCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subspaces",
		Short: "Parse things related to the x/subspaces module",
	}

	cmd.AddCommand(
		subspaceCmd(parseCfg),
		subspacesCmd(parseCfg),
	)

	return cmd
}
