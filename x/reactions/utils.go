package reactions

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/forbole/juno/v5/types/utils"

	juno "github.com/forbole/juno/v5/types"

	"github.com/forbole/juno/v5/node/remote"

	reactionstypes "github.com/desmos-labs/desmos/v6/x/reactions/types"

	"github.com/desmos-labs/athena/types"
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

// GetReactionIDFromEvent returns the reaction ID from the given event
func GetReactionIDFromEvent(event abci.Event) (uint32, error) {
	reactionIDAttr, err := utils.FindAttributeByKey(event, reactionstypes.AttributeKeyReactionID)
	if err != nil {
		return 0, err
	}

	reactionID, err := reactionstypes.ParseReactionID(reactionIDAttr.Value)
	if err != nil {
		return 0, err
	}

	return reactionID, nil
}

// GetRegisteredReactionIDFromEvent returns the registered reaction ID from the given event
func GetRegisteredReactionIDFromEvent(event abci.Event) (uint32, error) {
	registeredReactionIDAttr, err := utils.FindAttributeByKey(event, reactionstypes.AttributeKeyRegisteredReactionID)
	if err != nil {
		return 0, err
	}

	registeredReactionID, err := reactionstypes.ParseRegisteredReactionID(registeredReactionIDAttr.Value)
	if err != nil {
		return 0, err
	}

	return registeredReactionID, nil
}

// -------------------------------------------------------------------------------------------------------------------

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
