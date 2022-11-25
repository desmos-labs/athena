package posts

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/desmos-labs/djuno/v2/utils"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v4/node/remote"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshPostsData refreshes all the posts' data for the given subspace
func (m *Module) RefreshPostsData(height int64, subspaceID uint64) error {
	postsTxs, err := m.queryPostsTxs(height, subspaceID)
	if err != nil {
		return err
	}

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

		postTxs := m.getPostTxHashes(postsTxs, post)

		// Refresh transactions
		for _, tx := range postTxs {
			err = m.db.SavePostTx(tx)
			if err != nil {
				return fmt.Errorf("error while saving post transaction: %s", tx.Hash)
			}
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
			posts = append(posts, types.NewPost(post, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return posts, nil
}

// queryPostsTxs queries all the posts transactions made inside the given subspace
func (m *Module) queryPostsTxs(height int64, subspaceID uint64) ([]*coretypes.ResultTx, error) {
	queries := []string{
		// MsgCreatePost
		fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
			poststypes.EventTypeCreatePost,
			poststypes.AttributeKeySubspaceID,
			subspaceID,
			height,
		),
		// MsgEditPost
		fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
			poststypes.EventTypeEditPost,
			poststypes.AttributeKeySubspaceID,
			subspaceID,
			height,
		),
		// MsgAddPostAttachment
		fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
			poststypes.EventTypeAddPostAttachment,
			poststypes.AttributeKeySubspaceID,
			subspaceID,
			height,
		),
		// MsgRemovePostAttachment
		fmt.Sprintf("%s.%s='%d' AND tx.height <= %d",
			poststypes.EventTypeRemovePostAttachment,
			poststypes.AttributeKeySubspaceID,
			subspaceID,
			height,
		),
	}

	var txs []*coretypes.ResultTx
	for _, qry := range queries {
		resultTxs, err := utils.QueryTxs(m.node, qry)
		if err != nil {
			return nil, err
		}
		txs = append(txs, resultTxs...)
	}

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	return txs, nil
}

// filters the given transactions and returns only the ones related to the given post
func (m *Module) getPostTxHashes(txs []*coretypes.ResultTx, post types.Post) []types.PostTransaction {
	var postTxs []types.PostTransaction
	for _, tx := range txs {
		for _, event := range tx.TxResult.Events {
			if utils.HasSubspaceIDAndPostIDAttributes(event, post.SubspaceID, post.ID) {
				postTxs = append(postTxs, types.NewPostTransaction(
					post.SubspaceID,
					post.ID,
					hex.EncodeToString(tx.Hash),
				))
			}
		}
	}
	return postTxs
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
