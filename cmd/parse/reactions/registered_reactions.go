package profiles

import (
	"fmt"

	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/reactions"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

// registeredReactionsCmd returns a Cobra command that allows to refresh all the registered reactions
func registeredReactionsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "registered-reactions",
		Short: "Fetch all the posts reactions from the node and save them properly",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			subspacesModule := subspaces.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)
			reactionsModule := reactions.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Get the subspaces
			subspaces, err := subspacesModule.QueryAllSubspaces(height)
			if err != nil {
				return err
			}

			log.Info().Int64("height", height).Msg("refreshing registered reactions")
			for _, subspace := range subspaces {
				// Refresh the registered reactions
				err := reactionsModule.RefreshRegisteredReactionsData(height, subspace.ID)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
