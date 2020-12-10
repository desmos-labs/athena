package relationships

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	relationshipstypes "github.com/desmos-labs/desmos/x/relationships/types"
	desmosdb "github.com/desmos-labs/djuno/database"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to properly handle relationships-related messages
func HandleMsg(tx *juno.Tx, msg sdk.Msg, db *desmosdb.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {

	// Users
	case *relationshipstypes.MsgCreateRelationship:
		return handleMsgCreateRelationship(desmosMsg, db)

	case *relationshipstypes.MsgDeleteRelationship:
		return HandleMsgDeleteRelationship(desmosMsg, db)

	case *relationshipstypes.MsgBlockUser:
		return HandleMsgBlockUser(desmosMsg, db)

	case *relationshipstypes.MsgUnblockUser:
		return HandleMsgUnblockUser(desmosMsg, db)
	}

	return nil
}

// handleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func handleMsgCreateRelationship(msg *relationshipstypes.MsgCreateRelationship, db *desmosdb.DesmosDb) error {
	return db.SaveRelationship(msg.Sender, msg.Receiver, msg.Subspace)
}

// HandleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func HandleMsgDeleteRelationship(msg *relationshipstypes.MsgDeleteRelationship, db *desmosdb.DesmosDb) error {
	return db.DeleteRelationship(msg.User, msg.Counterparty, msg.Subspace)
}

// HandleMsgBlockUser allows to handle a MsgBlockUser properly
func HandleMsgBlockUser(msg *relationshipstypes.MsgBlockUser, db *desmosdb.DesmosDb) error {
	return db.SaveBlockage(msg.Blocker, msg.Blocked, msg.Reason, msg.Subspace)
}

// HandleMsgUnblockUser allows to handle a MsgUnblockUser properly
func HandleMsgUnblockUser(msg *relationshipstypes.MsgUnblockUser, db *desmosdb.DesmosDb) error {
	return db.RemoveBlockage(msg.Blocker, msg.Blocked, msg.Subspace)
}
