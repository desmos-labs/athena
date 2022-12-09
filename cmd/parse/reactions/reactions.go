package reactions

import (
	"fmt"

	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"

	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/reactions"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

// reactionsCmd returns a Cobra command that allows to refresh all the reactions
func reactionsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reactions [[subspace-id]]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Refresh all the reactions data",
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
			postsModule := posts.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)
			reactionsModule := reactions.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Get the subspaces
			log.Info().Int64("height", height).Msg("refreshing reactions")

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
				// Get the posts
				posts, err := postsModule.QuerySubspacePosts(height, subspaceID)
				if err != nil {
					return err
				}

				for _, post := range posts {
					// Refresh the reactions
					err = reactionsModule.RefreshReactionsData(height, post.SubspaceID, post.ID)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
