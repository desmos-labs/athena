package profilesscore

import (
	"github.com/desmos-labs/athena/types"
	"github.com/desmos-labs/athena/x/profiles"
)

type Database interface {
	profiles.Database
	SaveApplicationLinkScore(score *types.ProfileScore) error
}
