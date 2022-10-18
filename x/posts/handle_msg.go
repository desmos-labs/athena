package posts

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"

	"github.com/desmos-labs/djuno/v2/x/filters"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 || !filters.ShouldMsgBeParsed(msg) {
		return nil
	}

	switch desmosMsg := msg.(type) {
	case *poststypes.MsgCreatePost:
		return m.handleMsgCreatePost(tx, index, desmosMsg)

	case *poststypes.MsgEditPost:
		return m.handleMsgEditPost(tx, desmosMsg)

	case *poststypes.MsgDeletePost:
		return m.handleMsgDeletePost(tx, desmosMsg)

	case *poststypes.MsgAddPostAttachment:
		return m.handleMsgAddPostAttachment(tx, desmosMsg)

	case *poststypes.MsgRemovePostAttachment:
		return m.handleMsgRemovePostAttachment(tx, desmosMsg)

	case *poststypes.MsgAnswerPoll:
		return m.handleMsgAnswerPoll(tx, desmosMsg)
	}

	log.Debug().Str("module", "reactions").Str("message", proto.MessageName(msg)).
		Int64("height", tx.Height).Msg("handled message")

	return nil
}

// handleMsgCreatePost handles a MsgCreatePost
func (m *Module) handleMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) error {
	event, err := tx.FindEventByType(index, poststypes.EventTypeCreatePost)
	if err != nil {
		return err
	}
	postIDStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return err
	}
	postID, err := poststypes.ParsePostID(postIDStr)
	if err != nil {
		return err
	}

	// Update the post
	err = m.updatePost(tx.Height, msg.SubspaceID, postID)
	if err != nil {
		return err
	}

	// Update the post attachments
	return m.updatePostAttachments(tx.Height, msg.SubspaceID, postID)
}

// handleMsgEditPost handles a MsgEditPost
func (m *Module) handleMsgEditPost(tx *juno.Tx, msg *poststypes.MsgEditPost) error {
	return m.updatePost(tx.Height, msg.SubspaceID, msg.PostID)
}

// handleMsgDeletePost handles a MsgDeletePost
func (m *Module) handleMsgDeletePost(tx *juno.Tx, msg *poststypes.MsgDeletePost) error {
	return m.db.DeletePost(tx.Height, msg.SubspaceID, msg.PostID)
}

// handleMsgAddPostAttachment handles a MsgAddPostAttachment
func (m *Module) handleMsgAddPostAttachment(tx *juno.Tx, msg *poststypes.MsgAddPostAttachment) error {
	return m.updatePostAttachments(tx.Height, msg.SubspaceID, msg.PostID)
}

// handleMsgRemovePostAttachment handles a MsgRemovePostAttachment
func (m *Module) handleMsgRemovePostAttachment(tx *juno.Tx, msg *poststypes.MsgRemovePostAttachment) error {
	return m.db.DeletePostAttachment(tx.Height, msg.SubspaceID, msg.PostID, msg.AttachmentID)
}

// handleMsgAnswerPoll handles a MsgAnswerPoll
func (m *Module) handleMsgAnswerPoll(tx *juno.Tx, msg *poststypes.MsgAnswerPoll) error {
	answer := poststypes.NewUserAnswer(msg.SubspaceID, msg.PostID, msg.PollID, msg.AnswersIndexes, msg.Signer)
	return m.db.SavePollAnswer(types.NewPollAnswer(answer, tx.Height))
}
