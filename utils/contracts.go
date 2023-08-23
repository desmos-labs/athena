package utils

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// IsMsgSendTip tries reading the given MsgExecuteContract as if
// it's containing a send_tip message, and returns the inner message
func IsMsgSendTip(msg *wasmtypes.MsgExecuteContract) (*types.MsgSendTip, bool) {
	var msgTip types.TipMsg
	err := json.Unmarshal(msg.Msg, &msgTip)
	if err != nil {
		return nil, false
	}
	return msgTip.SendTip, true
}
