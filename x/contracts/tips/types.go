package tips

import (
	"bytes"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ContractQuery struct {
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

// --------------------------------------------------------------------------------------------------------------------

type TipMsg struct {
	SendTip *MsgSendTip `json:"send_tip,omitempty"`
}

type MsgSendTip struct {
	Amount sdk.Coins         `json:"amount"`
	Target *MsgSendTipTarget `json:"target"`
}

type MsgSendTipTarget struct {
	User    *TargetUser    `json:"user_target"`
	Content *TargetContent `json:"content_target"`
}

func (t *MsgSendTipTarget) Equal(u *MsgSendTipTarget) bool {
	tBz, _ := json.Marshal(t)
	uBz, _ := json.Marshal(u)
	return bytes.EqualFold(tBz, uBz)
}

type TargetUser struct {
	Receiver string `json:"receiver"`
}

type TargetContent struct {
	PostID string `json:"post_id"`
}
