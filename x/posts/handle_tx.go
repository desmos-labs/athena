package posts

import (
	abci "github.com/cometbft/cometbft/abci/types"
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/utils/events"
	"github.com/desmos-labs/athena/utils/transactions"

	"github.com/desmos-labs/athena/types"
)

// HandleTx handles the transaction events
func (m *Module) HandleTx(tx *juno.Tx) error {
	return transactions.ParseTxEvents(tx, map[string]func(tx *juno.Tx, event abci.Event) error{
		poststypes.EventTypeCreatePost:           m.parseCreatePostEvent,
		poststypes.EventTypeEditPost:             m.parseEditPostEvent,
		poststypes.EventTypeDeletePost:           m.parseDeletePostEvent,
		poststypes.EventTypeAddPostAttachment:    m.parseAddPostAttachmentEvent,
		poststypes.EventTypeRemovePostAttachment: m.parseRemovePostAttachmentEvent,
		poststypes.EventTypeAnswerPoll:           m.parseAnswerPollEvent,
	})
}

// -------------------------------------------------------------------------------------------------------------------

// parseCreatePostEvent handles the creation of a new post
func (m *Module) parseCreatePostEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the post
	err = m.updatePost(tx.Height, subspaceID, postID)
	if err != nil {
		return err
	}

	// Update the post attachments
	err = m.updatePostAttachments(tx.Height, subspaceID, postID)
	if err != nil {
		return err
	}

	// Save the related transaction
	return m.db.SavePostTx(types.NewPostTransaction(subspaceID, postID, tx.TxHash))
}

// parseEditPostEvent handles the edition of an existing post
func (m *Module) parseEditPostEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the post
	err = m.updatePost(tx.Height, subspaceID, postID)
	if err != nil {
		return err
	}

	// Save the related transaction
	return m.db.SavePostTx(types.NewPostTransaction(subspaceID, postID, tx.TxHash))
}

// parseDeletePostEvent handles the deletion of an existing post
func (m *Module) parseDeletePostEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.db.DeletePost(tx.Height, subspaceID, postID)
}

// parseAddPostReactionEvent handles the addition of a reaction to an existing post
func (m *Module) parseAddPostAttachmentEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	// Update the attachments
	err = m.updatePostAttachments(tx.Height, subspaceID, postID)
	if err != nil {
		return err
	}

	// Store the related post transaction
	return m.db.SavePostTx(types.NewPostTransaction(subspaceID, postID, tx.TxHash))
}

// parseRemovePostAttachmentEvent handles the removal of a reaction from an existing post
func (m *Module) parseRemovePostAttachmentEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	attachmentID, err := GetAttachmentIDFromEvent(event)
	if err != nil {
		return err
	}

	// Delete the attachment
	err = m.db.DeletePostAttachment(tx.Height, subspaceID, postID, attachmentID)
	if err != nil {
		return err
	}

	// Store the related post transaction
	return m.db.SavePostTx(types.NewPostTransaction(subspaceID, postID, tx.TxHash))
}

// parseAnswerPollEvent handles the answer to a poll
func (m *Module) parseAnswerPollEvent(tx *juno.Tx, event abci.Event) error {
	subspaceID, err := events.GetSubspaceIDFromEvent(event)
	if err != nil {
		return err
	}

	postID, err := GetPostIDFromEvent(event)
	if err != nil {
		return err
	}

	pollID, err := GetPollIDFromEvent(event)
	if err != nil {
		return err
	}

	return m.updatePollAnswers(tx.Height, subspaceID, postID, pollID)
}
