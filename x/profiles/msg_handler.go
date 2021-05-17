package profiles

import (
	"time"

	"github.com/desmos-labs/djuno/x/profiles/types"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	juno "github.com/desmos-labs/juno/types"

	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsg allows to handle different messages types for the profiles module
func HandleMsg(tx *juno.Tx, index int, msg sdk.Msg, db *desmosdb.DesmosDb) error {
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
		return HandleMsgDeleteRelationship(tx, desmosMsg, db)

	case *profilestypes.MsgBlockUser:
		return HandleMsgBlockUser(tx, desmosMsg, db)

	case *profilestypes.MsgUnblockUser:
		return HandleMsgUnblockUser(tx, desmosMsg, db)
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------

// handleMsgSaveProfile handles a MsgCreateProfile and properly stores the new profile inside the database
func handleMsgSaveProfile(tx *juno.Tx, index int, msg *profilestypes.MsgSaveProfile, db *desmosdb.DesmosDb) error {
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
func handleMsgDeleteProfile(tx *juno.Tx, msg *profilestypes.MsgDeleteProfile, db *desmosdb.DesmosDb) error {
	return db.DeleteProfile(msg.Creator, tx.Height)
}

// -----------------------------------------------------------------------------------------------------

// handleMsgRequestDTagTransfer handles a MsgRequestDTagTransfer storing the request into the database
func handleMsgRequestDTagTransfer(
	tx *juno.Tx, index int, msg *profilestypes.MsgRequestDTagTransfer, db *desmosdb.DesmosDb,
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
func handleMsgAcceptDTagTransfer(tx *juno.Tx, msg *profilestypes.MsgAcceptDTagTransfer, db *desmosdb.DesmosDb) error {
	return db.TransferDTag(types.NewDTagTransferRequestAcceptance(
		types.NewDTagTransferRequest(
			profilestypes.NewDTagTransferRequest(msg.NewDTag, msg.Sender, msg.Receiver),
			tx.Height,
		),
		msg.NewDTag,
	))
}

// handleDTagTransferRequestDeletion allows to delete an existing transfer request
func handleDTagTransferRequestDeletion(height int64, sender, receiver string, db *desmosdb.DesmosDb) error {
	return db.DeleteDTagTransferRequest(types.NewDTagTransferRequest(
		profilestypes.NewDTagTransferRequest("", sender, receiver),
		height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func handleMsgCreateRelationship(tx *juno.Tx, msg *profilestypes.MsgCreateRelationship, db *desmosdb.DesmosDb) error {
	return db.SaveRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.Sender, msg.Receiver, msg.Subspace),
		tx.Height,
	))
}

// HandleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func HandleMsgDeleteRelationship(tx *juno.Tx, msg *profilestypes.MsgDeleteRelationship, db *desmosdb.DesmosDb) error {
	return db.DeleteRelationship(types.NewRelationship(
		profilestypes.NewRelationship(msg.User, msg.Counterparty, msg.Subspace),
		tx.Height,
	))
}

// HandleMsgBlockUser allows to handle a MsgBlockUser properly
func HandleMsgBlockUser(tx *juno.Tx, msg *profilestypes.MsgBlockUser, db *desmosdb.DesmosDb) error {
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

// HandleMsgUnblockUser allows to handle a MsgUnblockUser properly
func HandleMsgUnblockUser(tx *juno.Tx, msg *profilestypes.MsgUnblockUser, db *desmosdb.DesmosDb) error {
	return db.RemoveBlockage(types.NewBlockage(
		profilestypes.NewUserBlock(msg.Blocker, msg.Blocked, "", msg.Subspace),
		tx.Height,
	))
}
