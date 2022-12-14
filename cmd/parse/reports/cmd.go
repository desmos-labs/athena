package reports

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewReportsCmd returns the Cobra command that allows to parse things related to the x/reports module
func NewReportsCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reports",
		Short: "Parse things related to the x/reports module",
	}

	cmd.AddCommand(
		reasonsCmd(parseCfg),
		reportsCmd(parseCfg),
	)

	return cmd
}
