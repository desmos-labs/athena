package profiles

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/juno/client"

	"github.com/desmos-labs/djuno/types"
	"github.com/desmos-labs/djuno/x/profiles/ibc"
	profilesutils "github.com/desmos-labs/djuno/x/profiles/utils"

	"github.com/cosmos/cosmos-sdk/codec"
	channeltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v2/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsg allows to handle different messages types for the profiles module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *profilestypes.MsgSaveProfile:
		return handleMsgSaveProfile(tx, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgDeleteProfile:
		return handleMsgDeleteProfile(tx, desmosMsg, db)

	case *profilestypes.MsgRequestDTagTransfer:
		return handleMsgRequestDTagTransfer(tx, index, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgAcceptDTagTransferRequest:
		return handleMsgAcceptDTagTransfer(tx, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgCancelDTagTransferRequest:
		return handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver, db)

	case *profilestypes.MsgRefuseDTagTransferRequest:
		return handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver, db)

	case *profilestypes.MsgCreateRelationship:
		return handleMsgCreateRelationship(tx, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgDeleteRelationship:
		return handleMsgDeleteRelationship(tx, desmosMsg, db)

	case *profilestypes.MsgBlockUser:
		return handleMsgBlockUser(tx, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgUnblockUser:
		return handleMsgUnblockUser(tx, desmosMsg, db)

	case *profilestypes.MsgLinkChainAccount:
		return handleMsgChainLink(tx, index, desmosMsg, profilesClient, cdc, db)

	case *profilestypes.MsgUnlinkChainAccount:
		return handleMsgUnlinkChainAccount(desmosMsg, db)

	case *profilestypes.MsgLinkApplication:
		return handleMsgLinkApplication(tx, desmosMsg, profilesClient, cdc, db)

	case *channeltypes.MsgRecvPacket:
		return ibc.HandlePacket(tx.Height, desmosMsg.Packet, profilesClient, cdc, db)

	case *channeltypes.MsgAcknowledgement:
		return ibc.HandlePacket(tx.Height, desmosMsg.Packet, profilesClient, cdc, db)

	case *channeltypes.MsgTimeout:
		return ibc.HandlePacket(tx.Height, desmosMsg.Packet, profilesClient, cdc, db)

	case *profilestypes.MsgUnlinkApplication:
		return handleMsgUnlinkApplication(desmosMsg, db)
	}

	log.Info().Str("module", "profiles").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// handleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func handleMsgSaveProfile(
	tx *juno.Tx, msg *profilestypes.MsgSaveProfile,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	addresses := []string{msg.Creator}
	return profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
}

// handleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func handleMsgDeleteProfile(tx *juno.Tx, msg *profilestypes.MsgDeleteProfile, db *desmosdb.Db) error {
	return db.DeleteProfile(msg.Creator, tx.Height)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgRequestDTagTransfer handles a MsgRequestDTagTransfer storing the request into the database
func handleMsgRequestDTagTransfer(
	tx *juno.Tx, index int, msg *profilestypes.MsgRequestDTagTransfer,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
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

	return db.SaveDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest(dTagToTrade, msg.Sender, msg.Receiver),
		tx.Height,
	))
}

// handleMsgAcceptDTagTransfer handles a MsgAcceptDTagTransfer effectively transferring the DTag
func handleMsgAcceptDTagTransfer(
	tx *juno.Tx, msg *profilestypes.MsgAcceptDTagTransferRequest,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	return profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
}

// handleDTagTransferRequestDeletion allows to delete an existing transfer request
func handleDTagTransferRequestDeletion(height int64, sender, receiver string, db *desmosdb.Db) error {
	return db.DeleteDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest("", sender, receiver),
		height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func handleMsgCreateRelationship(
	tx *juno.Tx, msg *profilestypes.MsgCreateRelationship,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Receiver, msg.Sender}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	return db.SaveRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.Sender, msg.Receiver, msg.Subspace),
		tx.Height,
	))
}

// handleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func handleMsgDeleteRelationship(tx *juno.Tx, msg *profilestypes.MsgDeleteRelationship, db *desmosdb.Db) error {
	return db.DeleteRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.User, msg.Counterparty, msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgBlockUser allows to handle a MsgBlockUser properly
func handleMsgBlockUser(
	tx *juno.Tx, msg *profilestypes.MsgBlockUser,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Blocked, msg.Blocker}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	return db.SaveBlockage(types.NewBlockage(
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
func handleMsgUnblockUser(tx *juno.Tx, msg *profilestypes.MsgUnblockUser, db *desmosdb.Db) error {
	return db.RemoveBlockage(types.NewBlockage(
		profilestypes.NewUserBlock(msg.Blocker, msg.Blocked, "", msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgChainLink allows to handle a MsgLinkChainAccount properly
func handleMsgChainLink(
	tx *juno.Tx, index int, msg *profilestypes.MsgLinkChainAccount,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.Signer}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
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
	err = cdc.UnpackAny(msg.ChainAddress, &address)
	if err != nil {
		return fmt.Errorf("error while unpacking address data: %s", err)
	}

	return db.SaveChainLink(types.NewChainLink(
		profilestypes.NewChainLink(msg.Signer, address, msg.Proof, msg.ChainConfig, creationTime),
		tx.Height,
	))
}

// handleMsgUnlinkChainAccount allows to handle a MsgUnlinkChainAccount properly
func handleMsgUnlinkChainAccount(msg *profilestypes.MsgUnlinkChainAccount, db *desmosdb.Db) error {
	return db.DeleteChainLink(msg.Owner, msg.Target, msg.ChainName)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgLinkApplication allows to handle a MsgLinkApplication properly
func handleMsgLinkApplication(
	tx *juno.Tx, msg *profilestypes.MsgLinkApplication,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *desmosdb.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.Sender}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating profiles: %s", strings.Join(addresses, ","))
	}

	res, err := profilesClient.UserApplicationLink(
		context.Background(),
		profilestypes.NewQueryUserApplicationLinkRequest(msg.Sender, msg.LinkData.Application, msg.LinkData.Username),
		client.GetHeightRequestHeader(tx.Height),
	)
	if err != nil {
		return fmt.Errorf("error while getting application link: %s", err)
	}

	return db.SaveApplicationLink(types.NewApplicationLink(res.Link, tx.Height))
}

// handleMsgUnlinkApplication allows to handle a MsgUnlinkApplication properly
func handleMsgUnlinkApplication(msg *profilestypes.MsgUnlinkApplication, db *desmosdb.Db) error {
	return db.DeleteApplicationLink(msg.Signer, msg.Application, msg.Username)
}
