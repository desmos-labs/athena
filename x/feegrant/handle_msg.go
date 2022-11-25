package feegrant

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	juno "github.com/forbole/juno/v4/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *feegrant.MsgGrantAllowance:
		return m.handleMsgGrantAllowance(tx, desmosMsg)
	case *feegrant.MsgRevokeAllowance:
		return m.handleMsgRevokeAllowance(tx, desmosMsg)
	}

	return nil
}

// handleMsgGrantAllowance handles a single MsgGrantAllowance message
func (m *Module) handleMsgGrantAllowance(tx *juno.Tx, msg *feegrant.MsgGrantAllowance) error {
	// Unpack interfaces
	err := msg.UnpackInterfaces(m.cdc)
	if err != nil {
		return fmt.Errorf("error while unpacking MsgGrantAllowance interfaces: %s", err)
	}

	// Get the allowance
	allowance, err := msg.GetFeeAllowanceI()
	if err != nil {
		return err
	}

	return m.db.SaveFeeGrant(types.NewFeeGrant(msg.Granter, msg.Grantee, allowance, tx.Height))
}

// handleMsgRevokeAllowance handles a single MsgRevokeAllowance message
func (m *Module) handleMsgRevokeAllowance(tx *juno.Tx, msg *feegrant.MsgRevokeAllowance) error {
	return m.db.DeleteFeeGrant(msg.Granter, msg.Grantee, tx.Height)
}
