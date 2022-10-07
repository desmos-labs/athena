package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

type FeeGrant struct {
	Granter   string
	Grantee   string
	Allowance feegrant.FeeAllowanceI
	Height    int64
}

func NewFeeGrant(granter string, grantee string, allowance feegrant.FeeAllowanceI, height int64) FeeGrant {
	return FeeGrant{
		Granter:   granter,
		Grantee:   grantee,
		Allowance: allowance,
		Height:    height,
	}
}

func (grant FeeGrant) GetSpendLimit() (sdk.Coins, error) {
	return getAllowanceSpendLimit(grant.Allowance)
}

func getAllowanceSpendLimit(allowance feegrant.FeeAllowanceI) (sdk.Coins, error) {
	switch grant := allowance.(type) {
	case *feegrant.BasicAllowance:
		return grant.SpendLimit, nil
	case *feegrant.PeriodicAllowance:
		return grant.Basic.SpendLimit, nil
	case *feegrant.AllowedMsgAllowance:
		allowance, err := grant.GetAllowance()
		if err != nil {
			return nil, err
		}
		return getAllowanceSpendLimit(allowance)
	default:
		return nil, fmt.Errorf("allowance type %T not supported", allowance)
	}
}

func (grant FeeGrant) GetExpirationDate() (*time.Time, error) {
	return getAllowanceExpirationDate(grant.Allowance)
}

func getAllowanceExpirationDate(allowance feegrant.FeeAllowanceI) (*time.Time, error) {
	switch grant := allowance.(type) {
	case *feegrant.BasicAllowance:
		return grant.Expiration, nil
	case *feegrant.PeriodicAllowance:
		return grant.Basic.Expiration, nil
	case *feegrant.AllowedMsgAllowance:
		allowance, err := grant.GetAllowance()
		if err != nil {
			return nil, err
		}
		return getAllowanceExpirationDate(allowance)
	default:
		return nil, fmt.Errorf("allowance type %T not supported", allowance)
	}
}
