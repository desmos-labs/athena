package builder

import (
	"github.com/forbole/juno/v3/types/config"

	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/domain"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/github"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/twitch"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/twitter"
)

func BuildModule(junoCfg config.Config, db profilesscore.Database) *profilesscore.Module {
	return profilesscore.NewModule([]profilesscore.Scorer{
		domain.NewScorer(),
		github.NewScorer(junoCfg),
		twitch.NewScorer(junoCfg),
		twitter.NewScorer(junoCfg),
	}, db)
}
