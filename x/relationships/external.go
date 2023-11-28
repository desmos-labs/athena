package relationships

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/query"
	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"
	"github.com/forbole/juno/v5/node/remote"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/athena/types"
)

// RefreshRelationshipsData refreshes all the relationships data for the given subspace
func (m *Module) RefreshRelationshipsData(height int64, subspaceID uint64) error {
	relationships, err := m.queryAllRelationships(height, subspaceID)
	if err != nil {
		return fmt.Errorf("error while getting relationships from gRPC: %s", err)
	}

	err = m.db.DeleteAllRelationships(height, subspaceID)
	if err != nil {
		return fmt.Errorf("error while deleting relationships: %s", err)
	}

	for _, relationship := range relationships {
		log.Debug().Uint64("subspace", relationship.SubspaceID).Str("creator", relationship.Creator).
			Str("counterparty", relationship.Counterparty).Msg("refreshing relationship")

		err = m.db.SaveRelationship(relationship)
		if err != nil {
			return fmt.Errorf("error while saving relationship: %s", err)
		}
	}

	return nil
}

// queryAllRelationships queries all the relationships for the given subspace from the node
func (m *Module) queryAllRelationships(height int64, subspaceID uint64) ([]types.Relationship, error) {
	var relationships []types.Relationship

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Relationships(
			remote.GetHeightRequestContext(context.Background(), height),
			&relationshipstypes.QueryRelationshipsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, relationship := range res.Relationships {
			relationships = append(relationships, types.NewRelationship(relationship, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return relationships, nil

}

// --------------------------------------------------------------------------------------------------------------------

// RefreshUserBlocksData refreshes all the user blocks data for the given subspace
func (m *Module) RefreshUserBlocksData(height int64, subspaceID uint64) error {
	userBlocks, err := m.queryAllUserBlocks(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllUserBlocks(height, subspaceID)
	if err != nil {
		return err
	}

	for _, block := range userBlocks {
		log.Debug().Uint64("subspace", block.SubspaceID).Str("blocker", block.Blocker).
			Str("blocked", block.Blocked).Msg("refreshing block")

		err = m.db.SaveUserBlock(block)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllUserBlocks queries all the user blocks for the given subspace from the node
func (m *Module) queryAllUserBlocks(height int64, subspaceID uint64) ([]types.Blockage, error) {
	var userBlocks []types.Blockage

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Blocks(
			remote.GetHeightRequestContext(context.Background(), height),
			&relationshipstypes.QueryBlocksRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, blockage := range res.Blocks {
			userBlocks = append(userBlocks, types.NewBlockage(blockage, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return userBlocks, nil

}
