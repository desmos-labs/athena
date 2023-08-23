package reactions

import (
	"context"

	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/juno/v5/node/remote"

	reactionstypes "github.com/desmos-labs/desmos/v6/x/reactions/types"

	"github.com/desmos-labs/djuno/v2/types"
)

func (m *Module) GetReactionID(tx *juno.Tx, index int) (uint32, error) {
	event, err := tx.FindEventByType(index, reactionstypes.EventTypeAddReaction)
	if err != nil {
		return 0, err
	}
	reactionIDStr, err := tx.FindAttributeByKey(event, reactionstypes.AttributeKeyReactionID)
	if err != nil {
		return 0, err
	}
	return reactionstypes.ParseReactionID(reactionIDStr)
}

func (m *Module) GetReaction(height int64, subspaceID uint64, postID uint64, reactionID uint32) (types.Reaction, error) {
	res, err := m.client.Reaction(
		remote.GetHeightRequestContext(context.Background(), height),
		&reactionstypes.QueryReactionRequest{
			SubspaceId: subspaceID,
			PostId:     postID,
			ReactionId: reactionID,
		},
	)
	if err != nil {
		return types.Reaction{}, err
	}

	return types.NewReaction(res.Reaction, height), nil
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
