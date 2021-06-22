package profiles

import (
	"fmt"
	"time"

	"github.com/desmos-labs/djuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/modules/messages"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	juno "github.com/desmos-labs/juno/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsg allows to handle different messages types for the profiles module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg,
	getAccounts messages.MessageAddressesParser, cdc codec.Marshaler, db *desmosdb.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *profilestypes.MsgSaveProfile:
		return handleMsgSaveProfile(tx, index, desmosMsg, db)

	case *profilestypes.MsgDeleteProfile:
		return handleMsgDeleteProfile(tx, desmosMsg, db)

	case *profilestypes.MsgRequestDTagTransfer:
		return handleMsgRequestDTagTransfer(tx, index, desmosMsg, db)

	case *profilestypes.MsgAcceptDTagTransfer:
		return handleMsgAcceptDTagTransfer(tx, desmosMsg, db)

	case *profilestypes.MsgCancelDTagTransfer:
		return handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver, db)

	case *profilestypes.MsgRefuseDTagTransfer:
		return handleDTagTransferRequestDeletion(tx.Height, desmosMsg.Sender, desmosMsg.Receiver, db)

	case *profilestypes.MsgCreateRelationship:
		return handleMsgCreateRelationship(tx, desmosMsg, db)

	case *profilestypes.MsgDeleteRelationship:
		return handleMsgDeleteRelationship(tx, desmosMsg, db)

	case *profilestypes.MsgBlockUser:
		return handleMsgBlockUser(tx, desmosMsg, db)

	case *profilestypes.MsgUnblockUser:
		return handleMsgUnblockUser(tx, desmosMsg, db)

	case *profilestypes.MsgLinkChainAccount:
		return handleMsgChainLink(tx, index, desmosMsg, cdc, db)

	case *profilestypes.MsgUnlinkChainAccount:
		return handleMsgUnlinkChainAccount(desmosMsg, db)
	}

	return saveAccounts(tx.Height, msg, getAccounts, cdc, db)
}

// -------------------------------------------------------------------------------------------------------------------

// saveAccounts stores the accounts included in the given messages inside the profile table
func saveAccounts(
	height int64, msg sdk.Msg, getAccounts messages.MessageAddressesParser, cdc codec.Marshaler, db *desmosdb.Db,
) error {
	accounts, err := getAccounts(cdc, msg)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = db.SaveUserIfNotExisting(account, height)
		if err != nil {
			return err
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// handleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func handleMsgSaveProfile(tx *juno.Tx, index int, msg *profilestypes.MsgSaveProfile, db *desmosdb.Db) error {
	event, err := tx.FindEventByType(index, profilestypes.EventTypeProfileSaved)
	if err != nil {
		return err
	}

	// Get creation date
	creationDateStr, err := tx.FindAttributeByKey(event, profilestypes.AttributeProfileCreationTime)
	if err != nil {
		return err
	}
	creationDate, err := time.Parse(time.RFC3339, creationDateStr)
	if err != nil {
		return err
	}

	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return err
	}

	newProfile, err := profilestypes.NewProfile(
		msg.DTag,
		msg.Nickname,
		msg.Bio,
		profilestypes.NewPictures(msg.ProfilePicture, msg.CoverPicture),
		creationDate,
		authtypes.NewBaseAccountWithAddress(address),
	)
	if err != nil {
		return err
	}

	return db.SaveProfile(types.NewProfile(newProfile, tx.Height))
}

// handleMsgDeleteProfile handles a MsgDeleteProfile correctly deleting the account present inside the database
func handleMsgDeleteProfile(tx *juno.Tx, msg *profilestypes.MsgDeleteProfile, db *desmosdb.Db) error {
	return db.DeleteProfile(msg.Creator, tx.Height)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgRequestDTagTransfer handles a MsgRequestDTagTransfer storing the request into the database
func handleMsgRequestDTagTransfer(
	tx *juno.Tx, index int, msg *profilestypes.MsgRequestDTagTransfer, db *desmosdb.Db,
) error {
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
func handleMsgAcceptDTagTransfer(tx *juno.Tx, msg *profilestypes.MsgAcceptDTagTransfer, db *desmosdb.Db) error {
	return db.TransferDTag(types.NewDTagTransferRequestAcceptance(
		types.NewDTagTransferRequest(
			profilestypes.NewDTagTransferRequest(msg.NewDTag, msg.Sender, msg.Receiver),
			tx.Height,
		),
		msg.NewDTag,
	))
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
func handleMsgCreateRelationship(tx *juno.Tx, msg *profilestypes.MsgCreateRelationship, db *desmosdb.Db) error {
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

// handleMsgBlockUser allows to handle a MsgBlockUser properly
func handleMsgBlockUser(tx *juno.Tx, msg *profilestypes.MsgBlockUser, db *desmosdb.Db) error {
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

func handleMsgChainLink(tx *juno.Tx, index int, msg *profilestypes.MsgLinkChainAccount, cdc codec.Marshaler, db *desmosdb.Db) error {
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

func handleMsgUnlinkChainAccount(msg *profilestypes.MsgUnlinkChainAccount, db *desmosdb.Db) error {
	return db.DeleteChainLink(msg.Owner, msg.Target, msg.ChainName)
}
