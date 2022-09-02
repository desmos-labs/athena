package posts

import (
	"context"
	"fmt"

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

	return types.NewPost(res.Post, height), nil
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
