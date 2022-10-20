package utils

import (
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	juno "github.com/forbole/juno/v3/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func HasSubspaceIDAndPostIDAttributes(event abci.Event, subspaceID uint64, postID uint64) bool {
	subspaceIDAttr, err := juno.FindAttributeByKey(event, poststypes.AttributeKeySubspaceID)
	if err != nil {
		return false
	}
	subspaceIDValue, err := subspacestypes.ParseSubspaceID(string(subspaceIDAttr.Value))
	if err != nil {
		return false
	}

	postIDAttr, err := juno.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return false
	}
	postIDValue, err := poststypes.ParsePostID(string(postIDAttr.Value))
	if err != nil {
		return false
	}

	return subspaceIDValue == subspaceID && postIDValue == postID
}
