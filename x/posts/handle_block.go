package posts

import (
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	juno "github.com/forbole/juno/v5/types"

	poststypes "github.com/desmos-labs/desmos/v7/x/posts/types"
	subspacestypes "github.com/desmos-labs/desmos/v7/x/subspaces/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(block *coretypes.ResultBlock, results *coretypes.ResultBlockResults, _ []*juno.Tx, _ *coretypes.ResultValidators) error {
	for _, event := range juno.FindEventsByType(results.EndBlockEvents, poststypes.EventTypeTalliedPoll) {
		// Get the subspace id
		subspaceIDStr, err := juno.FindAttributeByKey(event, subspacestypes.AttributeKeySubspaceID)
		if err != nil {
			return err
		}
		subspaceID, err := subspacestypes.ParseSubspaceID(subspaceIDStr.Value)
		if err != nil {
			return err
		}

		// Get the post id
		postIDStr, err := juno.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
		if err != nil {
			return err
		}
		postID, err := poststypes.ParsePostID(postIDStr.Value)
		if err != nil {
			return err
		}

		// Update the post
		err = m.updatePost(block.Block.Height, subspaceID, postID)
		if err != nil {
			return err
		}
	}

	return nil
}
