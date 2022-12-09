package types

import (
	"bytes"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ContractTypeTips = "tips"
)

type Contract struct {
	Address  string
	Type     string
	ConfigBz []byte
	Height   int64
}

func NewContract(address string, contractType string, contactConfigBz []byte, height int64) Contract {
	return Contract{
		Address:  address,
		Type:     contractType,
		ConfigBz: contactConfigBz,
		Height:   height,
	}
}

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
