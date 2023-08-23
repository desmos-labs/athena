package utils

import (
	abci "github.com/cometbft/cometbft/abci/types"
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"
	juno "github.com/forbole/juno/v5/types"
)

func HasSubspaceIDAndPostIDAttributes(event abci.Event, subspaceID uint64, postID uint64) bool {
	subspaceIDAttr, err := juno.FindAttributeByKey(event, poststypes.AttributeKeySubspaceID)
	if err != nil {
		return false
	}
	subspaceIDValue, err := subspacestypes.ParseSubspaceID(subspaceIDAttr.Value)
	if err != nil {
		return false
	}

	postIDAttr, err := juno.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return false
	}
	postIDValue, err := poststypes.ParsePostID(postIDAttr.Value)
	if err != nil {
		return false
	}

	return subspaceIDValue == subspaceID && postIDValue == postID
}
