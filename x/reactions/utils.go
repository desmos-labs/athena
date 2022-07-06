package reactions

import (
	"context"

	"github.com/forbole/juno/v3/node/remote"

	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// updateReaction updates the stored data about the given reaction at the specified height
func (m *Module) updateReaction(height int64, subspaceID uint64, postID uint64, reactionID uint32) error {
	res, err := m.client.Reaction(
		remote.GetHeightRequestContext(context.Background(), height),
		&reactionstypes.QueryReactionRequest{
			SubspaceId: subspaceID,
			PostId:     postID,
			ReactionId: reactionID,
		},
	)
	if err != nil {
		return err
	}

	return m.db.SaveReaction(types.NewReaction(res.Reaction, height))
}

// updateRegisteredReaction updates the stored data about the given registered reaction at the specified height
func (m *Module) updateRegisteredReaction(height int64, subspaceID uint64, reactionID uint32) error {
	// Get the registered reaction
	res, err := m.client.RegisteredReaction(
		remote.GetHeightRequestContext(context.Background(), height),
		&reactionstypes.QueryRegisteredReactionRequest{SubspaceId: subspaceID, ReactionId: reactionID},
	)
	if err != nil {
		return err
	}

	// Save the registered reaction
	return m.db.SaveRegisteredReaction(types.NewRegisteredReaction(res.RegisteredReaction, height))
}

// updateReactionParams updates the stored data about the given reaction params at the specified height
func (m *Module) updateReactionParams(height int64, subspaceID uint64) error {
	// Get the params
	res, err := m.client.ReactionsParams(
		remote.GetHeightRequestContext(context.Background(), height),
		&reactionstypes.QueryReactionsParamsRequest{SubspaceId: subspaceID},
	)
	if err != nil {
		return err
	}

	// Save the params
	return m.db.SaveReactionParams(types.NewReactionParams(res.Params, height))
}
