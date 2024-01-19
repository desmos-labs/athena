package posts

import (
	"context"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/types/utils"

	"github.com/forbole/juno/v5/node/remote"

	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"

	"github.com/desmos-labs/athena/types"
)

// GetPostIDFromEvent returns the post ID from the given event
func GetPostIDFromEvent(event abci.Event) (uint64, error) {
	attribute, err := utils.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return 0, err
	}
	return poststypes.ParsePostID(attribute.Value)
}

// GetAttachmentIDFromEvent returns the attachment ID from the given event
func GetAttachmentIDFromEvent(event abci.Event) (uint32, error) {
	attribute, err := utils.FindAttributeByKey(event, poststypes.AttributeKeyAttachmentID)
	if err != nil {
		return 0, err
	}
	return poststypes.ParseAttachmentID(attribute.Value)
}

// GetPollIDFromEvent returns the poll ID from the given event
func GetPollIDFromEvent(event abci.Event) (uint32, error) {
	attribute, err := utils.FindAttributeByKey(event, poststypes.AttributeKeyPollID)
	if err != nil {
		return 0, err
	}
	return poststypes.ParseAttachmentID(attribute.Value)
}

// -------------------------------------------------------------------------------------------------------------------

// updatePost updates the stored data about the given post at the specified height
func (m *Module) updatePost(height int64, subspaceID uint64, postID uint64) error {
	post, err := m.GetPost(height, subspaceID, postID)
	if err != nil {
		return err
	}
	return m.db.SavePost(post)
}

// GetPost gets the given post from the chain
func (m *Module) GetPost(height int64, subspaceID uint64, postID uint64) (types.Post, error) {
	res, err := m.client.Post(
		remote.GetHeightRequestContext(context.Background(), height),
		&poststypes.QueryPostRequest{SubspaceId: subspaceID, PostId: postID},
	)
	if err != nil {
		return types.Post{}, err
	}

	return types.NewPost(res.Post, height), nil
}

// updatePostAttachments updates the stored attachments for the post having the given id
func (m *Module) updatePostAttachments(height int64, subspaceID uint64, postID uint64) error {
	attachments, err := m.GetPostAttachments(height, subspaceID, postID)
	if err != nil {
		return fmt.Errorf("error while getting post attachments: %s", err)
	}

	for _, attachment := range attachments {
		err = m.db.SavePostAttachment(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPostAttachments gets the attachments for the post having the given id
func (m *Module) GetPostAttachments(height int64, subspaceID uint64, postID uint64) ([]types.PostAttachment, error) {
	var attachments []types.PostAttachment
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.PostAttachments(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QueryPostAttachmentsRequest{
				SubspaceId: subspaceID,
				PostId:     postID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, attachment := range res.Attachments {
			attachments = append(attachments, types.NewPostAttachment(attachment, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return attachments, nil
}

// updatePollAnswers updates the stored answers for the poll having the given id
func (m *Module) updatePollAnswers(height int64, subspaceID uint64, postID uint64, pollID uint32) error {
	answers, err := m.GetPollAnswers(height, subspaceID, postID, pollID)
	if err != nil {
		return fmt.Errorf("error while getting poll answers: %s", err)
	}

	for _, answer := range answers {
		err = m.db.SavePollAnswer(answer)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPollAnswers gets the answers for the poll having the given id
func (m *Module) GetPollAnswers(height int64, subspaceID uint64, postID uint64, pollID uint32) ([]types.PollAnswer, error) {
	var answers []types.PollAnswer
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.PollAnswers(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QueryPollAnswersRequest{
				SubspaceId: subspaceID,
				PostId:     postID,
				PollId:     pollID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, answer := range res.Answers {
			answers = append(answers, types.NewPollAnswer(answer, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return answers, nil
}

// updateParams updates the stored params with the ones for the given height
func (m *Module) updateParams() error {
	height, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	// Get the params
	res, err := m.client.Params(
		remote.GetHeightRequestContext(context.Background(), height),
		&poststypes.QueryParamsRequest{},
	)
	if err != nil {
		return err
	}

	// Save the params
	return m.db.SavePostsParams(types.NewPostsParams(res.Params, height))
}
