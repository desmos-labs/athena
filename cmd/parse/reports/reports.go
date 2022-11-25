package reports

import (
	"fmt"

	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/node/remote"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/reports"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

// reportsCmd returns a Cobra command that allows to refresh all the reports
func reportsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "reports",
		Short: "Fetch all the reports from the node and save them properly",
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
			reportsModule := reports.NewModule(parseCtx.Node, grpcConnection, parseCtx.EncodingConfig.Marshaler, db)

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

			log.Info().Int64("height", height).Msg("refreshing reports")
			for _, subspace := range subspaces {
				err := reportsModule.RefreshReportsData(height, subspace.ID)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}
