package handlers

import (
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/notifications"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
	"time"
)

// HandleMsgCreatePost allows to properly handle the given msg present inside the specified tx at the specific
// index. It creates a new Post object from it, stores it inside the database and later sends out any
// push notification using Firebase Cloud Messaging.
func HandleMsgCreatePost(tx juno.Tx, index int, msg poststypes.MsgCreatePost, db database.DesmosDb) error {
	post, err := createAndStorePostFromMsgCreatePost(tx, index, msg, db)
	if err != nil {
		return err
	}

	return notifications.SendPostNotifications(*post, db)
}

// createAndStorePostFromMsgCreatePost allows to properly handle a MsgCreatePostEvent by storing inside the
// database the post that has been created with such message.
// After the post has been saved, it is returned for other uses.
func createAndStorePostFromMsgCreatePost(
	tx juno.Tx, index int, msg poststypes.MsgCreatePost, db database.DesmosDb,
) (*poststypes.Post, error) {
	// Get the post id
	event, err := tx.FindEventByType(index, poststypes.EventTypePostCreated)
	if err != nil {
		return nil, err
	}
	postIDStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return nil, err
	}
	postID, err := poststypes.ParsePostID(postIDStr)
	if err != nil {
		return nil, err
	}

	// Get the creation date
	createdString, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostCreationTime)
	if err != nil {
		return nil, err
	}
	created, err := time.Parse(time.RFC3339, createdString)
	if err != nil {
		return nil, err
	}

	// Create the post
	post := poststypes.NewPost(postID, msg.ParentID, msg.Message, msg.AllowsComments,
		msg.Subspace, msg.OptionalData, created, msg.Creator)

	if msg.Attachments != nil {
		post = post.WithAttachments(msg.Attachments)
	}

	if msg.PollData != nil {
		post = post.WithPollData(*msg.PollData)
	}

	log.Info().
		Str("id", postID.String()).
		Str("owner", post.Creator.String()).
		Msg("saving post")

	// Save the post
	err = db.SavePost(post)
	if err != nil {
		return nil, err
	}

	return &post, err
}

// ____________________________________

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func HandleMsgEditPost(tx juno.Tx, index int, msg poststypes.MsgEditPost, db database.DesmosDb) error {
	// Get the edit date
	event, err := tx.FindEventByType(index, poststypes.EventTypePostCreated)
	if err != nil {
		return err
	}
	editDateString, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostEditTime)
	if err != nil {
		return err
	}
	editDate, err := time.Parse(time.RFC3339, editDateString)
	if err != nil {
		return err
	}

	return db.EditPost(msg.PostID, msg.Message, editDate)
}
