package profiles

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewProfilesCmd returns the Cobra command that allows to fix all the things related to the x/profiles module
func NewProfilesCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Parse things related to the x/profiles module",
	}

	cmd.AddCommand(
		profilesCmd(parseCfg),
		applicationLinksCmd(parseCfg),
		chainLinksCmd(parseCfg),
	)

	return cmd
}
