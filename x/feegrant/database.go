package feegrant

import (
	"time"

	"github.com/desmos-labs/athena/types"
)

type Database interface {
	SaveFeeGrant(grant types.FeeGrant) error
	DeleteFeeGrant(granter string, grantee string, height int64) error
	DeleteExpiredFeeGrants(time time.Time) error
}
