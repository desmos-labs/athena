package posts

import (
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	"github.com/gogo/protobuf/proto"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
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
		return m.handleMsgAddPostAttachment(tx, index, desmosMsg)

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

	return m.updatePost(tx.Height, msg.SubspaceID, postID)
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
func (m *Module) handleMsgAddPostAttachment(tx *juno.Tx, index int, msg *poststypes.MsgAddPostAttachment) error {
	event, err := tx.FindEventByType(index, poststypes.EventTypeAddPostAttachment)
	if err != nil {
		return err
	}
	attachmentIDStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return err
	}
	attachmentID, err := poststypes.ParseAttachmentID(attachmentIDStr)
	if err != nil {
		return err
	}

	attachment := poststypes.NewAttachment(msg.SubspaceID, msg.PostID, attachmentID, msg.Content.GetCachedValue().(poststypes.AttachmentContent))
	return m.db.SavePostAttachment(types.NewPostAttachment(attachment, tx.Height))
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
