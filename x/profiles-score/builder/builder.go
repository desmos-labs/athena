package builder

import (
	"github.com/forbole/juno/v5/types/config"

	"github.com/desmos-labs/athena/x/profiles-score/scorers/youtube"

	profilesscore "github.com/desmos-labs/athena/x/profiles-score"
	"github.com/desmos-labs/athena/x/profiles-score/scorers/domain"
	"github.com/desmos-labs/athena/x/profiles-score/scorers/github"
	"github.com/desmos-labs/athena/x/profiles-score/scorers/twitch"
	"github.com/desmos-labs/athena/x/profiles-score/scorers/twitter"
)

func BuildModule(junoCfg config.Config, db profilesscore.Database) *profilesscore.Module {
	return profilesscore.NewModule([]profilesscore.Scorer{
		domain.NewScorer(),
		github.NewScorer(junoCfg),
		twitch.NewScorer(junoCfg),
		twitter.NewScorer(junoCfg),
		youtube.NewScorer(junoCfg),
	}, db)
}
