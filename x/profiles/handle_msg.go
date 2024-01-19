package profiles

import (
	"github.com/desmos-labs/athena/x/filters"

	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *channeltypes.MsgRecvPacket:
		return m.handlePacket(tx, desmosMsg.Packet)

	case *channeltypes.MsgAcknowledgement:
		return m.handlePacket(tx, desmosMsg.Packet)

	case *channeltypes.MsgTimeout:
		return m.handlePacket(tx, desmosMsg.Packet)
	}

	log.Debug().Str("module", "profiles").Str("message", sdk.MsgTypeURL(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}
