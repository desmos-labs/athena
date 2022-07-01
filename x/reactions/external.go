package reactions

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/types/query"
	reactionstypes "github.com/desmos-labs/desmos/v4/x/reactions/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshRegisteredReactionsData refreshes the registered reactions data for the given subspace
func (m *Module) RefreshRegisteredReactionsData(height int64, subspaceID uint64) error {
	reactions, err := m.queryAllRegisteredReactions(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllRegisteredReactions(height, subspaceID)
	if err != nil {
		return err
	}

	for _, reaction := range reactions {
		log.Debug().Uint64("subspace", reaction.SubspaceID).Uint32("reaction", reaction.ID).Msg("refreshing registered reaction")

		err = m.db.SaveRegisteredReaction(reaction)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefreshReactionsData refreshes the reactions data for the given post
func (m *Module) RefreshReactionsData(height int64, subspaceID uint64, postID uint64) error {
	reactions, err := m.queryAllReactions(height, subspaceID, postID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllReactions(height, subspaceID, postID)
	if err != nil {
		return err
	}

	for _, reaction := range reactions {
		log.Debug().Uint64("subspace", reaction.SubspaceID).Uint32("reaction", reaction.ID).Msg("refreshing reaction")

		err = m.db.SaveReaction(reaction)
		if err != nil {
			return err
		}
	}

	return nil
}

// RefreshParamsData refreshes the reactions params for the given subspace
func (m *Module) RefreshParamsData(height int64, subspaceID uint64) error {
	log.Debug().Uint64("subspace", subspaceID).Msg("refreshing reactions params")
	return m.updateReactionParams(height, subspaceID)
}

// queryAllRegisteredReactions queries all the registered reactions for the given subspace from the node
func (m *Module) queryAllRegisteredReactions(height int64, subspaceID uint64) ([]types.RegisteredReaction, error) {
	var reactions []types.RegisteredReaction

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.RegisteredReactions(
			remote.GetHeightRequestContext(context.Background(), height),
			&reactionstypes.QueryRegisteredReactionsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, reaction := range res.RegisteredReactions {
			reactions = append(reactions, types.NewRegisteredReaction(reaction, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return reactions, nil
}

// queryAllReactions queries all the reactions for the given post from the node
func (m *Module) queryAllReactions(height int64, subspaceID uint64, postID uint64) ([]types.Reaction, error) {
	var reactions []types.Reaction

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Reactions(
			remote.GetHeightRequestContext(context.Background(), height),
			&reactionstypes.QueryReactionsRequest{
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

		for _, reaction := range res.Reactions {
			reactions = append(reactions, types.NewReaction(reaction, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return reactions, nil
}
