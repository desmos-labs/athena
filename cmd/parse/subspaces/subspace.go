package subspaces

import (
	"fmt"

	contractsbuilder "github.com/desmos-labs/djuno/v2/x/contracts/builder"

	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/node/remote"
	"github.com/forbole/juno/v5/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/profiles"
	"github.com/desmos-labs/djuno/v2/x/reactions"
	"github.com/desmos-labs/djuno/v2/x/relationships"
	"github.com/desmos-labs/djuno/v2/x/reports"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

// subspaceCmd returns a Cobra command that allows to refresh a single subspace and all the data within it
func subspaceCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "subspace [subspace-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Refresh all the data related to the given subspace",
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
			profilesModule := profiles.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			subspacesModule := subspaces.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			postsModule := posts.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			reactionsModule := reactions.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			relationshipsModule := relationships.NewModule(profilesModule, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			reportsModule := reports.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)
			contractsModule := contractsbuilder.BuildModule(config.Cfg, parseCtx.Node, grpcConnection, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Refresh x/subspace data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing subspace")
			err = subspacesModule.RefreshSubspaceData(height, subspaceID)
			if err != nil {
				return err
			}

			// Refresh x/relationships data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing relationships")
			err = relationshipsModule.RefreshRelationshipsData(height, subspaceID)
			if err != nil {
				return err
			}

			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing user blocks")
			err = relationshipsModule.RefreshUserBlocksData(height, subspaceID)
			if err != nil {
				return err
			}

			// Refresh x/posts data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing posts")
			err = postsModule.RefreshPostsData(height, subspaceID)
			if err != nil {
				return nil
			}

			// Refresh x/reactions data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing reactions params")
			err = reactionsModule.RefreshParamsData(height, subspaceID)
			if err != nil {
				return err
			}

			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing registered reactions")
			err = reactionsModule.RefreshRegisteredReactionsData(height, subspaceID)
			if err != nil {
				return nil
			}

			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing reactions")
			posts, err := postsModule.QuerySubspacePosts(height, subspaceID)
			if err != nil {
				return err
			}

			for _, post := range posts {
				err = reactionsModule.RefreshReactionsData(height, post.SubspaceID, post.ID)
				if err != nil {
					return err
				}
			}

			// Refresh x/reports data
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing reports reasons")
			err = reportsModule.RefreshReasonsData(height, subspaceID)
			if err != nil {
				return err
			}

			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing reports")
			err = reportsModule.RefreshReportsData(height, subspaceID)
			if err != nil {
				return err
			}

			// Refresh smart contracts details
			log.Info().Int64("height", height).Uint64("subspace id", subspaceID).Msg("refreshing contracts")
			err = contractsModule.RefreshData(height, subspaceID)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
