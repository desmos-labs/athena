package posts

import (
	"time"

	"github.com/desmos-labs/djuno/x/posts/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/notifications"
)

// MsgHandler allows to handle different message types from the posts module
func MsgHandler(tx *juno.Tx, index int, msg sdk.Msg, database *database.DesmosDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case *poststypes.MsgCreatePost:
		return handleMsgCreatePost(tx, index, desmosMsg, database)

	case *poststypes.MsgEditPost:
		return handleMsgEditPost(tx, index, desmosMsg, database)

	// Reactions
	case *poststypes.MsgRegisterReaction:
		return handleMsgRegisterReaction(tx, desmosMsg, database)

	case *poststypes.MsgAddPostReaction:
		return handleMsgAddPostReaction(tx, index, database)

	case *poststypes.MsgRemovePostReaction:
		return handleMsgRemovePostReaction(tx, index, database)

	// Polls
	case *poststypes.MsgAnswerPoll:
		return handleMsgAnswerPoll(tx, desmosMsg, database)
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgCreatePost allows to properly handle the given msg present inside the specified tx at the specific
// index. It creates a new Post object from it, stores it inside the database and later sends out any
// push notification using Firebase Cloud Messaging.
func handleMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost, db *database.DesmosDb) error {
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
	tx *juno.Tx, index int, msg *poststypes.MsgCreatePost, db *database.DesmosDb,
) (*poststypes.Post, error) {
	event, err := tx.FindEventByType(index, poststypes.EventTypePostCreated)
	if err != nil {
		return nil, err
	}

	// Get the post id
	postID, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return nil, err
	}

	// Get the creation time
	creationTimeStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostCreationTime)
	if err != nil {
		return nil, err
	}
	creationTime, err := time.Parse(time.RFC3339, creationTimeStr)
	if err != nil {
		return nil, err
	}

	// Create the post
	post := poststypes.NewPost(
		postID,
		msg.ParentID,
		msg.Message,
		msg.AllowsComments,
		msg.Subspace,
		msg.OptionalData,
		msg.Attachments,
		msg.PollData,
		creationTime,
		time.Time{},
		msg.Creator,
	)

	log.Info().Str("id", postID).Str("owner", post.Creator).Msg("saving post")

	// Save the post
	err = db.SavePost(types.NewPost(&post, tx.Height))
	if err != nil {
		return nil, err
	}

	return &post, err
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func handleMsgEditPost(tx *juno.Tx, index int, msg *poststypes.MsgEditPost, db *database.DesmosDb) error {
	event, err := tx.FindEventByType(index, poststypes.EventTypePostCreated)
	if err != nil {
		return err
	}

	editDateStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostEditTime)
	if err != nil {
		return err
	}
	editDate, err := time.Parse(time.RFC3339, editDateStr)
	if err != nil {
		return err
	}

	// Get the post
	post, err := db.GetPostByID(msg.PostID)
	if err != nil {
		return err
	}

	// Update the post
	post.Message = msg.Message
	post.LastEdited = editDate

	if msg.Attachments != nil {
		post.Attachments = msg.Attachments
	}

	if msg.PollData != nil {
		post.PollData = msg.PollData
	}

	return db.SavePost(types.NewPost(post, tx.Height))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func handleMsgAnswerPoll(tx *juno.Tx, msg *poststypes.MsgAnswerPoll, db *database.DesmosDb) error {
	return db.SaveUserPollAnswer(types.NewUserPollAnswer(
		msg.PostID,
		poststypes.NewUserAnswer(msg.UserAnswers, msg.Answerer),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// getReactionFromTxEvent creates a new PostReaction object from the event having the given type and associated
// to the message having the given inside the inside the given tx.
func getReactionFromTxEvent(tx *juno.Tx, index int, eventType string) (string, poststypes.PostReaction, error) {
	event, err := tx.FindEventByType(index, eventType)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	postID, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	user, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionOwner)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	value, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionValue)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	shortCode, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyReactionShortCode)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	return postID, poststypes.NewPostReaction(shortCode, value, user), nil
}

// HandleMsgAddPostReaction allows to properly handle the adding of a reaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func handleMsgAddPostReaction(tx *juno.Tx, index int, db *database.DesmosDb) error {
	postID, reaction, err := getReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionAdded)
	if err != nil {
		return err
	}

	err = db.SavePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
	if err != nil {
		return err
	}

	return notifications.SendReactionNotifications(postID, reaction, db)
}

// HandleMsgRemovePostReaction allows to properly handle the removal of a reaction from a post by
// deleting the specified reaction from the database.
func handleMsgRemovePostReaction(tx *juno.Tx, index int, db *database.DesmosDb) error {
	postID, reaction, err := getReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionRemoved)
	if err != nil {
		return err
	}

	return db.RemovePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgRegisterReaction handles a MsgRegisterReaction by storing the new reaction inside the database.
func handleMsgRegisterReaction(tx *juno.Tx, msg *poststypes.MsgRegisterReaction, db *database.DesmosDb) error {
	return db.RegisterReactionIfNotPresent(types.NewRegisteredReaction(
		poststypes.NewRegisteredReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace),
		tx.Height,
	))
}
