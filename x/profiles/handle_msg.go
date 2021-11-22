package profiles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forbole/juno/v2/node/remote"

	"github.com/gogo/protobuf/proto"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	juno "github.com/forbole/juno/v2/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
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

	case *profilestypes.MsgCreateRelationship:
		return m.handleMsgCreateRelationship(tx, desmosMsg)

	case *profilestypes.MsgDeleteRelationship:
		return m.handleMsgDeleteRelationship(tx, desmosMsg)

	case *profilestypes.MsgBlockUser:
		return m.handleMsgBlockUser(tx, desmosMsg)

	case *profilestypes.MsgUnblockUser:
		return m.handleMsgUnblockUser(tx, desmosMsg)

	case *profilestypes.MsgLinkChainAccount:
		return m.handleMsgChainLink(tx, index, desmosMsg)

	case *profilestypes.MsgUnlinkChainAccount:
		return m.handleMsgUnlinkChainAccount(desmosMsg)

	case *profilestypes.MsgLinkApplication:
		return m.handleMsgLinkApplication(tx, desmosMsg)

	case *channeltypes.MsgRecvPacket:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *channeltypes.MsgAcknowledgement:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *channeltypes.MsgTimeout:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *profilestypes.MsgUnlinkApplication:
		return m.handleMsgUnlinkApplication(desmosMsg)
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
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	err := m.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	event, err := tx.FindEventByType(index, profilestypes.EventTypeDTagTransferRequest)
	if err != nil {
		return err
	}

	dTagToTrade, err := tx.FindAttributeByKey(event, profilestypes.AttributeDTagToTrade)
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

// handleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func (m *Module) handleMsgCreateRelationship(tx *juno.Tx, msg *profilestypes.MsgCreateRelationship) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	err := m.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	return m.db.SaveRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.Sender, msg.Receiver, msg.Subspace),
		tx.Height,
	))
}

// handleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func (m *Module) handleMsgDeleteRelationship(tx *juno.Tx, msg *profilestypes.MsgDeleteRelationship) error {
	return m.db.DeleteRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.User, msg.Counterparty, msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgBlockUser allows to handle a MsgBlockUser properly
func (m *Module) handleMsgBlockUser(tx *juno.Tx, msg *profilestypes.MsgBlockUser) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Blocked, msg.Blocker}
	err := m.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	return m.db.SaveBlockage(types.NewBlockage(
		profilestypes.NewUserBlock(
			msg.Blocker,
			msg.Blocked,
			msg.Reason,
			msg.Subspace,
		),
		tx.Height,
	))
}

// handleMsgUnblockUser allows to handle a MsgUnblockUser properly
func (m *Module) handleMsgUnblockUser(tx *juno.Tx, msg *profilestypes.MsgUnblockUser) error {
	return m.db.RemoveBlockage(types.NewBlockage(
		profilestypes.NewUserBlock(msg.Blocker, msg.Blocked, "", msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgChainLink allows to handle a MsgLinkChainAccount properly
func (m *Module) handleMsgChainLink(tx *juno.Tx, index int, msg *profilestypes.MsgLinkChainAccount) error {
	// Update the involved account profile
	addresses := []string{msg.Signer}
	err := m.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	// Get the creation time
	event, err := tx.FindEventByType(index, profilestypes.EventTypeLinkChainAccount)
	if err != nil {
		return err
	}
	creationTimeStr, err := tx.FindAttributeByKey(event, profilestypes.AttributeChainLinkCreationTime)
	if err != nil {
		return err
	}
	creationTime, err := time.Parse(time.RFC3339, creationTimeStr)
	if err != nil {
		return err
	}

	// Unpack the address data
	var address profilestypes.AddressData
	err = m.cdc.UnpackAny(msg.ChainAddress, &address)
	if err != nil {
		return fmt.Errorf("error while unpacking address data: %s", err)
	}

	return m.db.SaveChainLink(types.NewChainLink(
		profilestypes.NewChainLink(msg.Signer, address, msg.Proof, msg.ChainConfig, creationTime),
		tx.Height,
	))
}

// handleMsgUnlinkChainAccount allows to handle a MsgUnlinkChainAccount properly
func (m *Module) handleMsgUnlinkChainAccount(msg *profilestypes.MsgUnlinkChainAccount) error {
	return m.db.DeleteChainLink(msg.Owner, msg.Target, msg.ChainName)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgLinkApplication allows to handle a MsgLinkApplication properly
func (m *Module) handleMsgLinkApplication(tx *juno.Tx, msg *profilestypes.MsgLinkApplication) error {
	// Update the involved account profile
	addresses := []string{msg.Sender}
	err := m.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	res, err := m.profilesClient.UserApplicationLink(
		context.Background(),
		profilestypes.NewQueryUserApplicationLinkRequest(msg.Sender, msg.LinkData.Application, msg.LinkData.Username),
		remote.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		return fmt.Errorf("error while getting application link: %s", err)
	}

	return m.db.SaveApplicationLink(types.NewApplicationLink(res.Link, tx.Height))
}

// handleMsgUnlinkApplication allows to handle a MsgUnlinkApplication properly
func (m *Module) handleMsgUnlinkApplication(msg *profilestypes.MsgUnlinkApplication) error {
	return m.db.DeleteApplicationLink(msg.Signer, msg.Application, msg.Username)
}
