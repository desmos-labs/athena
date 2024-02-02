package authz

import (
	"time"

	"github.com/desmos-labs/athena/v2/types"
)

type Database interface {
	SaveAuthzGrant(grant types.AuthzGrant) error
	DeleteAuthzGrant(granter string, grantee string, msgTypeURL string, height int64) error
	DeleteExpiredGrants(time time.Time) error
}
