package profiles

import (
	"fmt"

	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/profiles"
)

// applicationLinksCmd returns a Cobra command that allows to fix the application links for all the profiles
func applicationLinksCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "application-links",
		Short: "Fetch the application links stored on chain and save them",
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

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
			profilesModule := profiles.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)

			// Refresh the application links
			log.Info().Int64("height", height).Msg("refreshing applications links")
			return profilesModule.RefreshApplicationLinks(height)
		},
	}
}
