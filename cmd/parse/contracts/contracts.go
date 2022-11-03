package profiles

import (
	"fmt"

	contractsbuilder "github.com/desmos-labs/djuno/v2/x/contracts/builder"

	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/forbole/juno/v3/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
)

// contractsCmd returns a Cobra command that allows to refresh the smart contracts data for a single subspace
func contractsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "all [subspace-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Refresh all the smart contracts data related to the given subspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			subspaceID, err := subspacestypes.ParseSubspaceID(args[0])
			if err != nil {
				return err
			}

			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			remoteCfg, ok := config.Cfg.Node.Details.(*remote.Details)
			if !ok {
				panic(fmt.Errorf("cannot run DJuno on local node"))
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
			contractsModule := contractsbuilder.BuildModule(config.Cfg, parseCtx.Node, grpcConnection, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Refresh the smart contracts data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing contracts")
			err = contractsModule.RefreshData(height, subspaceID)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
