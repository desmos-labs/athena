package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/djuno/notifications"
	jtypes "github.com/desmos-labs/juno/types"
)

// GetReactionFromTxEvent creates a new PostReaction object from the event having the given type and associated
// to the message having the given inside the inside the given tx.
func GetReactionFromTxEvent(tx jtypes.Tx, index int, eventType string) (*posts.PostID, *posts.PostReaction, error) {
	event, err := FindEventByType(tx, index, eventType)
	if err != nil {
		return nil, nil, err
	}

	postIDStr, err := FindAttributeByKey(tx, event, "post_id")
	if err != nil {
		return nil, nil, err
	}
	postID, err := posts.ParsePostID(postIDStr)
	if err != nil {
		return nil, nil, err
	}

	userStr, err := FindAttributeByKey(tx, event, "reaction_user")
	if err != nil {
		return nil, nil, err
	}
	user, err := sdk.AccAddressFromBech32(userStr)
	if err != nil {
		return nil, nil, err
	}

	value, err := FindAttributeByKey(tx, event, "reaction_value")
	if err != nil {
		return nil, nil, err
	}

	shortCode, err := FindAttributeByKey(tx, event, "reaction_shortcode")
	if err != nil {
		return nil, nil, err
	}

	reaction := posts.NewPostReaction(shortCode, value, user)
	return &postID, &reaction, nil
}

// ____________________________________

// HandleMsgAddPostReaction allows to properly handle the adding of a reaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func HandleMsgAddPostReaction(tx jtypes.Tx, index int, db db.DesmosDb) error {
	postID, reaction, err := GetReactionFromTxEvent(tx, index, "post_reaction_added")
	if err != nil {
		return err
	}

	err = db.SaveReaction(*postID, reaction)
	if err != nil {
		return err
	}

	return notifications.SendReactionNotifications(*postID, reaction, db)
}

// ____________________________________

// HandleMsgRemovePostReaction allows to properly handle the removal of a reaction from a post by
// deleting the specified reaction from the database.
func HandleMsgRemovePostReaction(tx jtypes.Tx, index int, db db.DesmosDb) error {
	postID, reaction, err := GetReactionFromTxEvent(tx, index, "post_reaction_removed")
	if err != nil {
		return err
	}

	return db.RemoveReaction(*postID, reaction)
}

// ____________________________________

// HandleMsgRegisterReaction handles a MsgRegisterReaction by storing the new reaction inside the database.
func HandleMsgRegisterReaction(msg posts.MsgRegisterReaction, db db.DesmosDb) error {
	reaction := posts.NewReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace)
	_, err := db.RegisterReactionIfNotPresent(reaction)
	return err
}
