package tips

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type tipsContractQuery struct {
	Config *configQuery `json:"config,omitempty"`
}

type configQuery struct {
}

type configResponse struct {
	Admin      string          `json:"admin,omitempty"`
	SubspaceID string          `json:"subspace_id,omitempty"`
	Fees       *tipsFeesConfig `json:"service_fee,omitempty"`
}

type tipsFeesConfig struct {
	Percentage *percentageFee `json:"percentage,omitempty"`
	Fixed      *fixedFee      `json:"fixed,omitempty"`
}

type percentageFee struct {
	Value string `json:"value"`
}

type fixedFee struct {
	Amount sdk.Coins `json:"amount"`
}
