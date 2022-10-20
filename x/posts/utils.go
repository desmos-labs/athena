package posts

import (
	"context"
	"encoding/hex"
	"fmt"
	"sort"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/desmos-labs/djuno/v2/utils"

	"github.com/forbole/juno/v3/node/remote"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

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

	return m.getFullPostDetails(height, res.Post)
}

func (m *Module) getFullPostDetails(height int64, post poststypes.Post) (types.Post, error) {
	var txs []*coretypes.ResultTx

	msgCreatePostQuery := fmt.Sprintf("%s.%s='%d' AND %s.%s=%d AND tx.height <= %d",
		poststypes.EventTypeCreatePost,
		poststypes.AttributeKeySubspaceID,
		post.SubspaceID,
		poststypes.EventTypeCreatePost,
		poststypes.AttributeKeyPostID,
		post.ID,
		height,
	)
	msgCreatePostTxs, err := utils.QueryTxs(m.node, msgCreatePostQuery)
	if err != nil {
		return types.Post{}, err
	}
	txs = append(txs, msgCreatePostTxs...)

	msgEditPostsQuery := fmt.Sprintf("%s.%s='%d' AND %s.%s=%d AND tx.height <= %d",
		poststypes.EventTypeEditPost,
		poststypes.AttributeKeySubspaceID,
		post.SubspaceID,
		poststypes.EventTypeEditPost,
		poststypes.AttributeKeyPostID,
		post.ID,
		height,
	)
	msgEditPostsTxs, err := utils.QueryTxs(m.node, msgEditPostsQuery)
	if err != nil {
		return types.Post{}, err
	}
	txs = append(txs, msgEditPostsTxs...)

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	txHashes := make([]string, len(txs))
	for i, tx := range txs {
		txHashes[i] = hex.EncodeToString(tx.Tx.Hash())
	}
	return types.NewPost(post, txHashes, height), nil
}

// updatePostAttachments updates the stored attachments for the post having the given id
func (m *Module) updatePostAttachments(height int64, subspaceID uint64, postID uint64) error {
	attachments, err := m.getPostAttachments(height, subspaceID, postID)
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

func (m *Module) getPostAttachments(height int64, subspaceID uint64, postID uint64) ([]types.PostAttachment, error) {
	var attachments []types.PostAttachment
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.PostAttachments(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QueryPostAttachmentsRequest{
				SubspaceId: subspaceID,
				PostId:     postID,
				Pagination: nil,
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
