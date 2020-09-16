package handlers

import (
	relationshipstypes "github.com/desmos-labs/desmos/x/relationships/types"
	desmosdb "github.com/desmos-labs/djuno/database"
)

// HandleMsgCreateRelationship allows to handle a MsgCreateRelationship properly
func HandleMsgCreateRelationship(msg relationshipstypes.MsgCreateRelationship, db desmosdb.DesmosDb) error {
	return db.SaveRelationship(msg.Sender, msg.Receiver, msg.Subspace)
}

// HandleMsgDeleteRelationship allows to handle a MsgDeleteRelationship properly
func HandleMsgDeleteRelationship(msg relationshipstypes.MsgDeleteRelationship, db desmosdb.DesmosDb) error {
	return db.DeleteRelationship(msg.Sender, msg.Counterparty, msg.Subspace)
}

// HandleMsgBlockUser allows to handle a MsgBlockUser properly
func HandleMsgBlockUser(msg relationshipstypes.MsgBlockUser, db desmosdb.DesmosDb) error {
	return db.SaveBlockage(msg.Blocker, msg.Blocked, msg.Reason, msg.Subspace)
}

// HandleMsgUnblockUser allows to handle a MsgUnblockUser properly
func HandleMsgUnblockUser(msg relationshipstypes.MsgUnblockUser, db desmosdb.DesmosDb) error {
	return db.RemoveBlockage(msg.Blocker, msg.Blocked, msg.Subspace)
}
