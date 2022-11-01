package builder

import (
	"github.com/forbole/juno/v3/types/config"

	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/domain"

	"github.com/desmos-labs/djuno/v2/database"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/github"
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers/twitter"
)

func BuildModule(junoCfg config.Config, db *database.Db) *profilesscore.Module {
	return profilesscore.NewModule([]profilesscore.Scorer{
		twitter.NewScorer(junoCfg),
		github.NewScorer(junoCfg),
		domain.NewScorer(),
	}, db)
}
