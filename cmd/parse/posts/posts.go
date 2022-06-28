package profiles

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/forbole/juno/v3/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

// postsCmd returns a Cobra command that allows to refresh all the posts
func postsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Fetch all the posts and their data from the node and save them properly",
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

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Get the subspaces
			subs, err := subspacesModule.QueryAllSubspaces(height)
			if err != nil {
				return err
			}

			for _, subspace := range subs {
				log.Debug().Int64("height", height).Uint64("subspace", subspace.ID).Msg("refreshing posts")
				err = postsModule.RefreshPostsData(height, subspace.ID)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
