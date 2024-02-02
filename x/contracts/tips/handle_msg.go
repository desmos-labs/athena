package tips

import (
	"fmt"

	"github.com/desmos-labs/athena/v2/utils"

	"github.com/cosmos/cosmos-sdk/x/authz"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/v2/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessagesModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *wasmtypes.MsgInstantiateContract:
		return m.handleMsgInstantiateContract(tx, index)
	case *wasmtypes.MsgExecuteContract:
		return m.handleMsgExecuteContract(tx, desmosMsg)
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// handleMsgInstantiateContract handles a MsgInstantiateContract instance by refreshing the stored tips contracts
func (m *Module) handleMsgInstantiateContract(tx *juno.Tx, index int) error {
	// Refresh the configuration
	address, err := m.base.ParseContractAddress(tx, index)
	if err != nil {
		return err
	}

	if !m.cfg.IsContractSupported(address) {
		return nil
	}

	// Store the contract base data
	err = m.base.HandleMsgInstantiateContract(tx, index, types.ContractTypeTips)
	if err != nil {
		return err
	}

	return m.refreshContractConfig(tx.Height, address)
}

// --------------------------------------------------------------------------------------------------------------------

// handleMsgExecuteContract handles a MsgExecuteContract that contains a send_tip message by storing the tip details
func (m *Module) handleMsgExecuteContract(tx *juno.Tx, msg *wasmtypes.MsgExecuteContract) error {
	if !m.cfg.IsContractSupported(msg.Contract) {
		return nil
	}

	msgSendTip, ok := utils.IsMsgSendTip(msg)
	if !ok {
		return nil
	}

	// If it's not a send_tip message, then return
	if msgSendTip == nil {
		return nil
	}

	// Get all the send_tip messages to the same target from within the transaction
	msgSendTips, err := m.getSameTargetSendTipMessages(tx, msgSendTip)
	if err != nil {
		return err
	}

	// Combine all the send_tip messages into one
	combinedMsg := m.combineMsgSendTips(msgSendTips)

	// Get the contract configuration
	config, err := m.getContractConfig(tx.Height, msg.Contract)
	if err != nil {
		return fmt.Errorf("error while getting contract config: %s", err)
	}

	// Convert the data into a Tip object
	tip, err := m.convertMsgSendTip(msg.Sender, combinedMsg, config, tx.Height)
	if err != nil {
		return err
	}

	// Don't store the tip if the target is a post and the post doesn't exist
	if postTarget, ok := tip.Target.(types.PostTarget); ok {
		subspaceID, err := subspacestypes.ParseSubspaceID(config.SubspaceID)
		if err != nil {
			return err
		}

		found, err := m.db.HasPost(tx.Height, subspaceID, postTarget.PostID)
		if err != nil {
			return err
		}

		// Skip the storing of the tip
		if !found {
			return nil
		}
	}

	// Save the tip
	return m.db.SaveTip(tip)
}

// getSameTargetSendTipMessages iterates over the given transaction, and extracts the
// inner send_tip message out of all MsgExecuteContract instances
func (m *Module) getSameTargetSendTipMessages(tx *juno.Tx, msg *types.MsgSendTip) ([]*types.MsgSendTip, error) {
	msgExecContracts, err := m.extractMsgExecuteContracts(tx.GetMsgs())
	if err != nil {
		return nil, err
	}

	var msgs []*types.MsgSendTip
	for _, msgExecContract := range msgExecContracts {
		msgSendTip, ok := utils.IsMsgSendTip(msgExecContract)
		if !ok {
			return nil, nil
		}

		if msgSendTip.Target.Equal(msg.Target) {
			msgs = append(msgs, msgSendTip)
		}
	}
	return msgs, nil
}

// extractMsgExecuteContracts extracts all the MsgExecuteContract from the given slice,
// performing a recursive search inside authz.MsgExec if neeeded
func (m *Module) extractMsgExecuteContracts(msgs []sdk.Msg) ([]*wasmtypes.MsgExecuteContract, error) {
	var msgExecContracts []*wasmtypes.MsgExecuteContract
	for _, msg := range msgs {
		switch sdkMsg := msg.(type) {
		case *wasmtypes.MsgExecuteContract:
			msgExecContracts = append(msgExecContracts, sdkMsg)
		case *authz.MsgExec:
			innerSdkMsg, err := sdkMsg.GetMessages()
			if err != nil {
				return nil, err
			}
			innerMsgs, err := m.extractMsgExecuteContracts(innerSdkMsg)
			if err != nil {
				return nil, err
			}
			msgExecContracts = append(msgExecContracts, innerMsgs...)
		}
	}
	return msgExecContracts, nil
}

// combineMsgSendTips combines the given send_tip messages into a single send_tip message
// that has the amount equals to the sum of all amounts.
// NOTE. All send_tip messages should have the same target
func (m *Module) combineMsgSendTips(msgs []*types.MsgSendTip) *types.MsgSendTip {
	if len(msgs) == 0 {
		return nil
	}

	amount := sdk.NewCoins()
	for _, msg := range msgs {
		amount = amount.Add(msg.Amount...)
	}

	return &types.MsgSendTip{
		Amount: amount,
		Target: msgs[0].Target,
	}
}

// convertMsgSendTip converts the given data into a types.Tip instance
func (m *Module) convertMsgSendTip(sender string, msg *types.MsgSendTip, config *configResponse, height int64) (types.Tip, error) {
	subspaceID, err := subspacestypes.ParseSubspaceID(config.SubspaceID)
	if err != nil {
		return types.Tip{}, fmt.Errorf("error while parsing subsapce id: %s", err)
	}

	var target types.Target
	switch {
	case msg.Target.User != nil:
		target = types.NewUserTarget(msg.Target.User.Receiver)
	case msg.Target.Content != nil:
		postID, err := poststypes.ParsePostID(msg.Target.Content.PostID)
		if err != nil {
			return types.Tip{}, fmt.Errorf("error while parsing post id: %s", err)
		}
		target = types.NewPostTarget(postID)
	}

	return types.NewTip(subspaceID, sender, target, msg.Amount, height), nil
}
