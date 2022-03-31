package posts

import (
	"time"

	"github.com/desmos-labs/djuno/v2/types"
	"github.com/desmos-labs/djuno/v2/x/posts/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch desmosMsg := msg.(type) {
	// Posts
	case *poststypes.MsgCreatePost:
		return m.handleMsgCreatePost(tx, index, desmosMsg)

	case *poststypes.MsgEditPost:
		return m.handleMsgEditPost(tx, index, desmosMsg)

	// Reactions
	case *poststypes.MsgRegisterReaction:
		return m.handleMsgRegisterReaction(tx, desmosMsg)

	case *poststypes.MsgAddPostReaction:
		return m.handleMsgAddPostReaction(tx, index, desmosMsg)

	case *poststypes.MsgRemovePostReaction:
		return m.handleMsgRemovePostReaction(tx, index)

	// Polls
	case *poststypes.MsgAnswerPoll:
		return m.handleMsgAnswerPoll(tx, desmosMsg)

	// Reports
	case *poststypes.MsgReportPost:
		return m.handleMsgReport(tx, desmosMsg)
	}

	return nil
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgCreatePost allows to properly handle the given msg present inside the specified tx at the specific
// index. It creates a new Post object from it, stores it inside the database and later sends out any
// push notification using Firebase Cloud Messaging.
func (m *Module) handleMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) error {
	// Update the involved account profile
	addresses := []string{msg.Creator}
	err := m.profilesModule.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return err
	}

	post, err := utils.GetPostFromMsgCreatePost(tx, index, msg)
	if err != nil {
		return err
	}

	// Save the post
	return m.db.SavePost(post)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func (m *Module) handleMsgEditPost(tx *juno.Tx, index int, msg *poststypes.MsgEditPost) error {
	// Update the involved accounts profiles
	addresses := []string{msg.Editor}
	err := m.profilesModule.UpdateProfiles(tx.Height, addresses)
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
	post, err := m.db.GetPostByID(msg.PostID)
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

	return m.db.SavePost(post)
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func (m *Module) handleMsgAnswerPoll(tx *juno.Tx, msg *poststypes.MsgAnswerPoll) error {
	// Update the involved account profile
	addresses := []string{msg.Answerer}
	err := m.profilesModule.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return err
	}

	return m.db.SaveUserPollAnswer(types.NewUserPollAnswer(
		poststypes.NewUserAnswer(msg.PostID, msg.Answerer, msg.Answers),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgAddPostReaction allows to properly handle the adding of a reaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func (m *Module) handleMsgAddPostReaction(tx *juno.Tx, index int, msg *poststypes.MsgAddPostReaction) error {
	// Update the involved account profile
	addresses := []string{msg.User}
	err := m.profilesModule.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return err
	}

	postID, reaction, err := utils.GetReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionAdded)
	if err != nil {
		return err
	}

	return m.db.SavePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
}

// HandleMsgRemovePostReaction allows to properly handle the removal of a reaction from a post by
// deleting the specified reaction from the database.
func (m *Module) handleMsgRemovePostReaction(tx *juno.Tx, index int) error {
	postID, reaction, err := utils.GetReactionFromTxEvent(tx, index, poststypes.EventTypePostReactionRemoved)
	if err != nil {
		return err
	}

	return m.db.RemovePostReaction(types.NewPostReaction(postID, reaction, tx.Height))
}

// -----------------------------------------------------------------------------------------------------

// HandleMsgRegisterReaction handles a MsgRegisterReaction by storing the new reaction inside the database.
func (m *Module) handleMsgRegisterReaction(tx *juno.Tx, msg *poststypes.MsgRegisterReaction) error {
	// Update the involved account profile
	addresses := []string{msg.Creator}
	err := m.profilesModule.UpdateProfiles(tx.Height, addresses)
	if err != nil {
		return err
	}

	return m.db.RegisterReactionIfNotPresent(types.NewRegisteredReaction(
		poststypes.NewRegisteredReaction(msg.Creator, msg.ShortCode, msg.Value, msg.Subspace),
		tx.Height,
	))
}

// -----------------------------------------------------------------------------------------------------

// handleMsgReport allows to handle a MsgReportPost properly
func (m *Module) handleMsgReport(tx *juno.Tx, msg *poststypes.MsgReportPost) error {
	return m.db.SaveReport(types.NewReport(
		poststypes.NewReport(
			msg.PostID,
			msg.ReportType,
			msg.Message,
			msg.User,
		),
		tx.Height,
	))
}
