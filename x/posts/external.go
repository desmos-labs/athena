package posts

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	juno "github.com/forbole/juno/v3/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/desmos-labs/djuno/v2/utils"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshPostsData refreshes all the posts' data for the given subspace
func (m *Module) RefreshPostsData(height int64, subspaceID uint64) error {
	posts, err := m.QuerySubspacePosts(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllPosts(height, subspaceID)
	if err != nil {
		return fmt.Errorf("error while deleting subspace posts: %s", err)
	}

	// Refresh posts
	for _, post := range posts {
		log.Debug().Uint64("subspace", post.SubspaceID).Uint64("post", post.ID).Msg("refreshing post")

		err = m.db.SavePost(post)
		if err != nil {
			return fmt.Errorf("error while saving post: %s", err)
		}

		attachments, err := m.queryPostAttachments(height, post.SubspaceID, post.ID)
		if err != nil {
			return err
		}

		// Refresh attachments
		for _, attachment := range attachments {
			log.Debug().Uint64("subspace", attachment.SubspaceID).Uint64("post", attachment.PostID).
				Uint32("attachment", attachment.ID).Msg("refreshing attachment")

			err = m.db.SavePostAttachment(attachment)
			if err != nil {
				return fmt.Errorf("error while saving post attachment: %s", err)
			}

			// Refresh poll answers
			if _, isPoll := attachment.Content.GetCachedValue().(*poststypes.Poll); isPoll {
				log.Debug().Uint64("subspace", attachment.SubspaceID).Uint64("post", attachment.PostID).
					Uint32("poll", attachment.ID).Msg("refreshing poll answers")

				answers, err := m.queryPollAnswers(height, attachment.SubspaceID, attachment.PostID, attachment.ID)
				if err != nil {
					return err
				}

				for _, answer := range answers {
					err = m.db.SavePollAnswer(answer)
					if err != nil {
						return fmt.Errorf("errror while saving poll answer: %s", err)
					}
				}
			}
		}
	}

	return nil
}

// QuerySubspacePosts queries all the posts present inside the given subspace at the provided height
func (m *Module) QuerySubspacePosts(height int64, subspaceID uint64) ([]types.Post, error) {
	txs, err := m.queryPostsTxs(height, subspaceID)
	if err != nil {
		return nil, err
	}

	var posts []types.Post

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.SubspacePosts(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QuerySubspacePostsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, post := range res.Posts {
			txHashes := m.getPostTxHashes(txs, post)
			posts = append(posts, types.NewPost(post, txHashes, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return posts, nil
}

func (m *Module) queryPostsTxs(height int64, subspaceID uint64) ([]*coretypes.ResultTx, error) {
	var txs []*coretypes.ResultTx

	msgCreatePostQuery := fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
		poststypes.EventTypeCreatePost,
		poststypes.AttributeKeySubspaceID,
		subspaceID,
		height,
	)
	msgCreatePostTxs, err := utils.QueryTxs(m.node, msgCreatePostQuery)
	if err != nil {
		return nil, err
	}
	txs = append(txs, msgCreatePostTxs...)

	msgEditPostsQuery := fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
		poststypes.EventTypeEditPost,
		poststypes.AttributeKeySubspaceID,
		subspaceID,
		height,
	)
	msgEditPostsTxs, err := utils.QueryTxs(m.node, msgEditPostsQuery)
	if err != nil {
		return nil, err
	}
	txs = append(txs, msgEditPostsTxs...)

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	return txs, nil
}

func (m *Module) getPostTxHashes(txs []*coretypes.ResultTx, post poststypes.Post) []string {
	var txHashes []string
	for _, tx := range txs {
		if m.isCreatePostTx(tx, post.SubspaceID, post.ID) || m.isEditPostTx(tx, post.SubspaceID, post.ID) {
			txHashes = append(txHashes, hex.EncodeToString(tx.Tx.Hash()))
		}
	}
	return txHashes
}

func (m *Module) isCreatePostTx(tx *coretypes.ResultTx, subspaceID uint64, postID uint64) bool {
	event, err := juno.FindEventByType(tx.TxResult.Events, poststypes.EventTypeCreatePost)
	if err != nil {
		return false
	}
	return utils.HasSubspaceIDAndPostIDAttributes(event, subspaceID, postID)
}

func (m *Module) isEditPostTx(tx *coretypes.ResultTx, subspaceID uint64, postID uint64) bool {
	event, err := juno.FindEventByType(tx.TxResult.Events, poststypes.EventTypeEditPost)
	if err != nil {
		return false
	}
	return utils.HasSubspaceIDAndPostIDAttributes(event, subspaceID, postID)
}

// QuerySubspacePosts queries all the attachments for the given post at the provided height
func (m *Module) queryPostAttachments(height int64, subspaceID uint64, postID uint64) ([]types.PostAttachment, error) {
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

// QuerySubspacePosts queries all the posts present inside the given subspace at the provided height
func (m *Module) queryPollAnswers(height int64, subspaceID uint64, postID uint64, pollID uint32) ([]types.PollAnswer, error) {
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
