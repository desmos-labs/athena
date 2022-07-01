package posts

import (
	"context"
	"fmt"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/types"
)

// updatePost updates the stored data about the given post at the specified height
func (m *Module) updatePost(height int64, subspaceID uint64, postID uint64) error {
	// Get the post
	res, err := m.client.Post(
		remote.GetHeightRequestContext(context.Background(), height),
		&poststypes.QueryPostRequest{SubspaceId: subspaceID, PostId: postID},
	)
	if err != nil {
		return err
	}

	// Save the post
	return m.db.SavePost(types.NewPost(res.Post, height))
}

// updateParams updates the stored params with the ones for the given height
func (m *Module) updateParams(height int64) error {
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
