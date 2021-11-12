package posts

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	profilesutils "github.com/desmos-labs/djuno/x/profiles/utils"

	"github.com/desmos-labs/djuno/types"
	"github.com/desmos-labs/djuno/x/posts/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	juno "github.com/forbole/juno/v2/types"

	"github.com/desmos-labs/djuno/database"
)

// MsgHandler allows to handle different message types from the posts module
func MsgHandler(
	tx *juno.Tx, index int, msg sdk.Msg,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, database *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case *poststypes.MsgCreatePost:
		return handleMsgCreatePost(tx, index, desmosMsg, profilesClient, cdc, database)

	case *poststypes.MsgEditPost:
		return handleMsgEditPost(tx, index, desmosMsg, profilesClient, cdc, database)

	// Reactions
	case *poststypes.MsgRegisterReaction:
		return handleMsgRegisterReaction(tx, desmosMsg, profilesClient, cdc, database)

	case *poststypes.MsgAddPostReaction:
		return handleMsgAddPostReaction(tx, index, desmosMsg, profilesClient, cdc, database)

	case *poststypes.MsgRemovePostReaction:
		return handleMsgRemovePostReaction(tx, index, database)

	// Polls
	case *poststypes.MsgAnswerPoll:
		return handleMsgAnswerPoll(tx, desmosMsg, profilesClient, cdc, database)

	// Reports
	case *poststypes.MsgReportPost:
		return handleMsgReport(tx, desmosMsg, database)
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgCreatePost allows to properly handle the given msg present inside the specified tx at the specific
// index. It creates a new Post object from it, stores it inside the database and later sends out any
// push notification using Firebase Cloud Messaging.
func handleMsgCreatePost(
	tx *juno.Tx, index int, msg *poststypes.MsgCreatePost,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.Creator}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return err
	}

	post, err := utils.GetPostFromMsgCreatePost(tx, index, msg)
	if err != nil {
		return err
	}

	// Save the post
	return db.SavePost(post)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func handleMsgEditPost(
	tx *juno.Tx, index int, msg *poststypes.MsgEditPost,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Editor}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return err
	}

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

	if msg.Poll != nil {
		post.Poll = msg.Poll
	}

	return db.SavePost(post)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func handleMsgAnswerPoll(
	tx *juno.Tx, msg *poststypes.MsgAnswerPoll,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.Answerer}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return err
	}

	return db.SaveUserPollAnswer(types.NewUserPollAnswer(
		poststypes.NewUserAnswer(msg.PostID, msg.Answerer, msg.Answers),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAddPostReaction allows to properly handle the adding of a reaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func handleMsgAddPostReaction(
	tx *juno.Tx, index int, msg *poststypes.MsgAddPostReaction,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.User}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return err
	}

	postID, reaction, err := utils.GetReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionAdded)
	if err != nil {
		return err
	}

	return db.SavePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
}

// HandleMsgRemovePostReaction allows to properly handle the removal of a reaction from a post by
// deleting the specified reaction from the database.
func handleMsgRemovePostReaction(tx *juno.Tx, index int, db *database.Db) error {
	postID, reaction, err := utils.GetReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionRemoved)
	if err != nil {
		return err
	}

	return db.RemovePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgRegisterReaction handles a MsgRegisterReaction by storing the new reaction inside the database.
func handleMsgRegisterReaction(
	tx *juno.Tx, msg *poststypes.MsgRegisterReaction,
	profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) error {
	// Update the involved account profile
	addresses := []string{msg.Creator}
	err := profilesutils.UpdateProfiles(tx.Height, addresses, profilesClient, cdc, db)
	if err != nil {
		return err
	}

	return db.RegisterReactionIfNotPresent(types.NewRegisteredReaction(
		poststypes.NewRegisteredReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgReport allows to handle a MsgReportPost properly
func handleMsgReport(tx *juno.Tx, msg *poststypes.MsgReportPost, db *database.Db) error {
	return db.SaveReport(types.NewReport(
		poststypes.NewReport(
			msg.PostID,
			msg.ReportType,
			msg.Message,
			msg.User,
		),
		tx.Height,
	))
}
