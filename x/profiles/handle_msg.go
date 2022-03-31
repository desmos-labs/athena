package profiles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forbole/juno/v3/node/remote"

	"github.com/gogo/protobuf/proto"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	profilestypes "github.com/desmos-labs/desmos/v3/x/profiles/types"
	juno "github.com/forbole/juno/v3/types"

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

	case *profilestypes.MsgLinkChainAccount:
		return m.handleMsgChainLink(tx, index, desmosMsg)

	case *profilestypes.MsgUnlinkChainAccount:
		return m.handleMsgUnlinkChainAccount(tx.Height, desmosMsg)

	case *profilestypes.MsgLinkApplication:
		return m.handleMsgLinkApplication(tx, desmosMsg)

	case *channeltypes.MsgRecvPacket:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *channeltypes.MsgAcknowledgement:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *channeltypes.MsgTimeout:
		return m.handlePacket(tx.Height, desmosMsg.Packet)

	case *profilestypes.MsgUnlinkApplication:
		return m.handleMsgUnlinkApplication(tx.Height, desmosMsg)
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
	creationTimeStr, err := tx.FindAttributeByKey(event, profilestypes.AttributeKeyChainLinkCreationTime)
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
func (m *Module) handleMsgUnlinkChainAccount(height int64, msg *profilestypes.MsgUnlinkChainAccount) error {
	return m.db.DeleteChainLink(msg.Owner, msg.Target, msg.ChainName, height)
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

	res, err := m.profilesClient.ApplicationLinks(
		remote.GetHeightRequestContext(context.Background(), tx.Height),
		profilestypes.NewQueryApplicationLinksRequest(msg.Sender, msg.LinkData.Application, msg.LinkData.Username, nil),
	)
	if err != nil {
		return fmt.Errorf("error while getting application link: %s", err)
	}

	if len(res.Links) == 0 {
		return fmt.Errorf("no application link found on chain")
	}

	if len(res.Links) > 1 {
		return fmt.Errorf("duplicated application link found on chain")
	}

	return m.db.SaveApplicationLink(types.NewApplicationLink(res.Links[0], tx.Height))
}

// handleMsgUnlinkApplication allows to handle a MsgUnlinkApplication properly
func (m *Module) handleMsgUnlinkApplication(height int64, msg *profilestypes.MsgUnlinkApplication) error {
	return m.db.DeleteApplicationLink(msg.Signer, msg.Application, msg.Username, height)
}
