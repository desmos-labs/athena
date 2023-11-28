package parse

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"

	parseblocks "github.com/forbole/juno/v5/cmd/parse/blocks"
	parsegenesis "github.com/forbole/juno/v5/cmd/parse/genesis"

	parseauthz "github.com/desmos-labs/athena/cmd/parse/authz"
	parsecontracts "github.com/desmos-labs/athena/cmd/parse/contracts"
	parsefeegrant "github.com/desmos-labs/athena/cmd/parse/feegrant"
	parseposts "github.com/desmos-labs/athena/cmd/parse/posts"
	parseprofiles "github.com/desmos-labs/athena/cmd/parse/profiles"
	parsereactions "github.com/desmos-labs/athena/cmd/parse/reactions"
	parserelationships "github.com/desmos-labs/athena/cmd/parse/relationships"
	parsereports "github.com/desmos-labs/athena/cmd/parse/reports"
	parsesubspaces "github.com/desmos-labs/athena/cmd/parse/subspaces"
)

// NewParseCmd returns the Cobra command allowing to parse some data without having to re-sync the whole database
func NewParseCmd(parseCfg *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Apply some fixes without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parsecmdtypes.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		parsegenesis.NewGenesisCmd(parseCfg),
		parseblocks.NewBlocksCmd(parseCfg),
		parseauthz.NewAuthzCmd(parseCfg),
		parsecontracts.NewContractsCmd(parseCfg),
		parsefeegrant.NewFeeGrant(parseCfg),
		parseprofiles.NewProfilesCmd(parseCfg),
		parserelationships.NewRelationshipsCmd(parseCfg),
		parsesubspaces.NewSubspacesCmd(parseCfg),
		parseposts.NewPostsCmd(parseCfg),
		parsereactions.NewReactionsCmd(parseCfg),
		parsereports.NewReportsCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}
