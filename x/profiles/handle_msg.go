package profiles

import (
	"github.com/desmos-labs/djuno/v2/x/filters"

	"github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/gogo/protobuf/proto"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	juno "github.com/forbole/juno/v4/types"

	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *profilestypes.MsgSaveProfile:
		return m.handleMsgSaveProfile(tx, desmosMsg)

	case *profilestypes.MsgDeleteProfile:
		return m.handleMsgDeleteProfile(tx, desmosMsg)

	case *profilestypes.MsgRequestDTagTransfer:
		return m.handleMsgRequestDTagTransfer(tx, index, desmosMsg)

	case *profilestypes.MsgAcceptDTagTransferRequest:
		return m.handleMsgAcceptDTagTransfer(tx, desmosMsg)

	case *profilestypes.MsgCancelDTagTransferRequest:
		return m.handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver)

	case *profilestypes.MsgRefuseDTagTransferRequest:
		return m.handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver)

	case *profilestypes.MsgLinkChainAccount:
		return m.handleMsgChainLink(tx, desmosMsg)

	case *profilestypes.MsgUnlinkChainAccount:
		return m.handleMsgUnlinkChainAccount(tx, desmosMsg)

	case *profilestypes.MsgLinkApplication:
		return m.handleMsgLinkApplication(tx, desmosMsg)

	case *channeltypes.MsgRecvPacket:
		return m.handlePacket(tx, desmosMsg.Packet)

	case *channeltypes.MsgAcknowledgement:
		return m.handlePacket(tx, desmosMsg.Packet)

	case *channeltypes.MsgTimeout:
		return m.handlePacket(tx, desmosMsg.Packet)

	case *profilestypes.MsgUnlinkApplication:
		return m.handleMsgUnlinkApplication(tx, desmosMsg)
	}

	log.Debug().Str("module", "profiles").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// handleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func (m *Module) handleMsgSaveProfile(tx *juno.Tx, msg *profilestypes.MsgSaveProfile) error {
	addresses := []string{msg.Creator}
	return m.UpdateProfiles(tx.Height, addresses)
}

// handleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func (m *Module) handleMsgDeleteProfile(tx *juno.Tx, msg *profilestypes.MsgDeleteProfile) error {
	return m.db.DeleteProfile(msg.Creator, tx.Height)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgRequestDTagTransfer handles a MsgRequestDTagTransfer storing the request into the database
func (m *Module) handleMsgRequestDTagTransfer(tx *juno.Tx, index int, msg *profilestypes.MsgRequestDTagTransfer) error {
	event, err := tx.FindEventByType(index, profilestypes.EventTypeDTagTransferRequest)
	if err != nil {
		return err
	}

	dTagToTrade, err := tx.FindAttributeByKey(event, profilestypes.AttributeKeyDTagToTrade)
	if err != nil {
		return err
	}

	return m.db.SaveDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest(dTagToTrade, msg.Sender, msg.Receiver),
		tx.Height,
	))
}

// handleMsgAcceptDTagTransfer handles a MsgAcceptDTagTransfer effectively transferring the DTag
func (m *Module) handleMsgAcceptDTagTransfer(tx *juno.Tx, msg *profilestypes.MsgAcceptDTagTransferRequest) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	return m.UpdateProfiles(tx.Height, addresses)
}

// handleDTagTransferRequestDeletion allows to delete an existing transfer request
func (m *Module) handleDTagTransferRequestDeletion(height int64, sender, receiver string) error {
	return m.db.DeleteDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest("", sender, receiver),
		height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgChainLink allows to handle a MsgLinkChainAccount properly
func (m *Module) handleMsgChainLink(tx *juno.Tx, msg *profilestypes.MsgLinkChainAccount) error {
	// Save the chain links
	err := m.updateUserChainLinks(tx.Height, msg.Signer)
	if err != nil {
		return err
	}

	// Update the default chain links
	return m.updateUserDefaultChainLinks(tx.Height, msg.Signer)
}

// handleMsgUnlinkChainAccount allows to handle a MsgUnlinkChainAccount properly
func (m *Module) handleMsgUnlinkChainAccount(tx *juno.Tx, msg *profilestypes.MsgUnlinkChainAccount) error {
	err := m.db.DeleteChainLink(msg.Owner, msg.Target, msg.ChainName, tx.Height)
	if err != nil {
		return err
	}

	// Update the default chain links
	return m.updateUserDefaultChainLinks(tx.Height, msg.Owner)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgLinkApplication allows to handle a MsgLinkApplication properly
func (m *Module) handleMsgLinkApplication(tx *juno.Tx, msg *profilestypes.MsgLinkApplication) error {
	return m.updateUserApplicationLinks(tx.Height, msg.Sender)
}

// handleMsgUnlinkApplication allows to handle a MsgUnlinkApplication properly
func (m *Module) handleMsgUnlinkApplication(tx *juno.Tx, msg *profilestypes.MsgUnlinkApplication) error {
	return m.db.DeleteApplicationLink(msg.Signer, msg.Application, msg.Username, tx.Height)
}
