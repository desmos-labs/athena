package profiles

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/desmos-labs/djuno/v2/database"
	profilesscorebuilder "github.com/desmos-labs/djuno/v2/x/profiles-score/builder"
)

// applicationLinksScoresCmd returns a Cobra command that allows to fix the application links scores for all the profiles
func applicationLinksScoresCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "application-links",
		Short: "Fetch the application links stored on chain and save them",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Get the latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			// Refresh the application link scores
			profilesScoreModule := profilesscorebuilder.BuildModule(config.Cfg, db)
			log.Info().Int64("height", height).Msg("refreshing applications links scores")
			return profilesScoreModule.RefreshApplicationLinksScores()
		},
	}
}
