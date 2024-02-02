package profiles

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/node/remote"
	"github.com/forbole/juno/v5/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/athena/v2/database"
	"github.com/desmos-labs/athena/v2/x/profiles"
)

// chainLinksCmd returns a Cobra command that allows to fix the chain links for all the profiles
func chainLinksCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "chain-links",
		Short: "Fetch the chain links stored on chain and save them",
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

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
			profilesModule := profiles.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Codec, db)

			// Refresh the chain links
			log.Info().Int64("height", height).Msg("refreshing chain links")
			return profilesModule.RefreshChainLinks(height)
		},
	}
}
