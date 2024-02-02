package relationships

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	"github.com/forbole/juno/v5/types/utils"

	"github.com/desmos-labs/athena/types"
)

func GetSubspaceFromEvent(event abci.Event) (uint64, error) {
	subspaceAttr, err := utils.FindAttributeByKey(event, relationshipstypes.AttributeKeySubspace)
	if err != nil {
		return 0, err
	}

	subspaceID, err := subspacestypes.ParseSubspaceID(subspaceAttr.Value)
	if err != nil {
		return 0, err
	}

	return subspaceID, nil
}

// GetCreatorFromEvent returns the creator of the relationship from the given event
func GetCreatorFromEvent(event abci.Event) (string, error) {
	creatorAttr, err := utils.FindAttributeByKey(event, relationshipstypes.AttributeRelationshipCreator)
	if err != nil {
		return "", err
	}
	return creatorAttr.Value, nil
}

// GetCounterpartyFromEvent returns the counterparty of the relationship from the given event
func GetCounterpartyFromEvent(event abci.Event) (string, error) {
	counterpartyAttr, err := utils.FindAttributeByKey(event, relationshipstypes.AttributeRelationshipCounterparty)
	if err != nil {
		return "", err
	}
	return counterpartyAttr.Value, nil
}

// GetBlockerFromEvent returns the blocker of the user block from the given event
func GetBlockerFromEvent(event abci.Event) (string, error) {
	blockerAttr, err := utils.FindAttributeByKey(event, relationshipstypes.AttributeKeyUserBlockBlocker)
	if err != nil {
		return "", err
	}
	return blockerAttr.Value, nil
}

// GetBlockedFromEvent returns the blocked of the user block from the given event
func GetBlockedFromEvent(event abci.Event) (string, error) {
	blockedAttr, err := utils.FindAttributeByKey(event, relationshipstypes.AttributeKeyUserBlockBlocked)
	if err != nil {
		return "", err
	}
	return blockedAttr.Value, nil
}

// --------------------------------------------------------------------------------------------------------------------

func (m *Module) updateUserBlock(height int64, subspaceID uint64, blocker, blocked string) error {
	blocks, err := m.client.Blocks(context.Background(), &relationshipstypes.QueryBlocksRequest{
		SubspaceId: subspaceID,
		Blocker:    blocker,
		Blocked:    blocked,
	})
	if err != nil {
		return err
	}

	for _, block := range blocks.Blocks {
		err = m.db.SaveUserBlock(types.NewBlockage(block, height))
		if err != nil {
			return err
		}
	}

	return nil
}
