package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/authz"
)

type AuthzGrant struct {
	Granter       string
	Grantee       string
	Authorization authz.Authorization
	Expiration    time.Time
	Height        int64
}

func NewAuthzGrant(granter, grantee string, authorization authz.Authorization, expiration time.Time, height int64) AuthzGrant {
	return AuthzGrant{
		Granter:       granter,
		Grantee:       grantee,
		Authorization: authorization,
		Expiration:    expiration,
		Height:        height,
	}
}
