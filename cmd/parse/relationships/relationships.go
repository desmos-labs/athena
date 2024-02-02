package relationships

import (
	"fmt"

	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/node/remote"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/athena/v2/database"
	"github.com/desmos-labs/athena/v2/x/profiles"
	"github.com/desmos-labs/athena/v2/x/relationships"
	"github.com/desmos-labs/athena/v2/x/subspaces"
)

// relationshipsCmd returns a Cobra command that allows to fix the relationships for all the profiles
func relationshipsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "relationships [[subspace-id]]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Refresh all the relationships data",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			remoteCfg, ok := config.Cfg.Node.Details.(*remote.Details)
			if !ok {
				panic(fmt.Errorf("cannot run Athena on local node"))
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
			subspacesModule := subspaces.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			profilesModule := profiles.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			relationshipsModule := relationships.NewModule(profilesModule, grpcConnection, parseCtx.EncodingConfig.Codec, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Get the subspaces
			log.Info().Int64("height", height).Msg("refreshing relationships")

			var subspaceIDs []uint64
			if len(args) > 0 {
				subspaceID, err := subspacestypes.ParseSubspaceID(args[0])
				if err != nil {
					return err
				}
				subspaceIDs = []uint64{subspaceID}
			} else {
				subs, err := subspacesModule.QueryAllSubspaces(height)
				if err != nil {
					return err
				}

				subspaceIDs = make([]uint64, len(subs))
				for i, subspace := range subs {
					subspaceIDs[i] = subspace.ID
				}
			}

			for _, subspaceID := range subspaceIDs {
				err := relationshipsModule.RefreshRelationshipsData(height, subspaceID)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
