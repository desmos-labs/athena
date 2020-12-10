package posts

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/notifications"
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// MsgHandler allows to handle different message types from the posts module
func MsgHandler(tx *juno.Tx, index int, msg sdk.Msg, database *database.DesmosDb) error {
	if len(tx.Logs) == 0 {
		log.Info().
			Str("module", "posts").
			Str("tx_hash", tx.TxHash).Int("msg_index", index).
			Msg("skipping message as it was not successful")
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
		return handleMsgRegisterReaction(desmosMsg, database)

	case *poststypes.MsgAddPostReaction:
		return handleMsgAddPostReaction(tx, index, database)

	case *poststypes.MsgRemovePostReaction:
		return handleMsgRemovePostReaction(tx, index, database)

	// Polls
	case *poststypes.MsgAnswerPoll:
		return handleMsgAnswerPoll(desmosMsg, database)
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
	err = db.SavePost(post)
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

	return db.EditPost(msg.PostID, msg.Message, msg.Attachments, msg.PollData, editDate)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func handleMsgAnswerPoll(msg *poststypes.MsgAnswerPoll, db *database.DesmosDb) error {
	return db.SaveUserPollAnswer(msg.PostID, poststypes.NewUserAnswer(msg.UserAnswers, msg.Answerer))
}

// -----------------------------------------------------------------------------------------------------

// getReactionFromTxEvent creates a new PostReaction object from the event having the given type and associated
// to the message having the given inside the inside the given tx.
func getReactionFromTxEvent(tx *juno.Tx, index int, eventType string) (string, *poststypes.PostReaction, error) {
	event, err := tx.FindEventByType(index, eventType)
	if err != nil {
		return "", nil, err
	}

	postID, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return "", nil, err
	}

	user, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionOwner)
	if err != nil {
		return "", nil, err
	}

	value, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionValue)
	if err != nil {
		return "", nil, err
	}

	shortCode, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyReactionShortCode)
	if err != nil {
		return "", nil, err
	}

	reaction := poststypes.NewPostReaction(shortCode, value, user)
	return postID, &reaction, nil
}

// HandleMsgAddPostReaction allows to properly handle the adding of a reaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func handleMsgAddPostReaction(tx *juno.Tx, index int, db *database.DesmosDb) error {
	postID, reaction, err := getReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionAdded)
	if err != nil {
		return err
	}

	err = db.SaveReaction(postID, reaction)
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

	return db.RemoveReaction(postID, reaction)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgRegisterReaction handles a MsgRegisterReaction by storing the new reaction inside the database.
func handleMsgRegisterReaction(msg *poststypes.MsgRegisterReaction, db *database.DesmosDb) error {
	reaction := poststypes.NewRegisteredReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace)
	return db.RegisterReactionIfNotPresent(reaction)
}
