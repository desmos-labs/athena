package relationships

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewRelationshipsCmd returns the Cobra command that allows to fix all the things related to the x/profiles module
func NewRelationshipsCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relationships",
		Short: "Parse things related to the x/relationships module",
	}

	cmd.AddCommand(
		relationshipsCmd(parseCfg),
		userBlocksCmd(parseCfg),
	)

	return cmd
}
