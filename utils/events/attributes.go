package events

import (
	abci "github.com/cometbft/cometbft/abci/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	juno "github.com/forbole/juno/v5/types"
)

// GetSubspaceIDFromEvent returns the subspace ID from the given event
func GetSubspaceIDFromEvent(event abci.Event) (uint64, error) {
	attribute, err := juno.FindAttributeByKey(event, subspacestypes.AttributeKeySubspaceID)
	if err != nil {
		return 0, err
	}
	return subspacestypes.ParseSubspaceID(attribute.Value)
}
