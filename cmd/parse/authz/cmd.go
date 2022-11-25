package authz

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewAuthzCmd returns the Cobra command that allows to parse things related to the x/authz module
func NewAuthzCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "authz",
		Short: "Parse things related to the x/authz module",
	}

	cmd.AddCommand(
		authorizationsCmd(parseCfg),
	)

	return cmd
}
