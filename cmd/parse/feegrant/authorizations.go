package feegrant

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/x/feegrant"

	"github.com/desmos-labs/djuno/v2/database"
)

// allowancesCmd returns a Cobra command that allows to refresh all the allowances
func allowancesCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "allowances",
		Short: "Fetch all the authorizations from the node and save them properly",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the module
			feegrantModule := feegrant.NewModule(parseCtx.Node, parseCtx.EncodingConfig.Codec, db)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Refresh the authorizations
			log.Info().Int64("height", height).Msg("refreshing fee grant allowances")
			return feegrantModule.RefreshFeeGrants(height)
		},
	}
}
