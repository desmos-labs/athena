package types

import (
	relationshipstypes "github.com/desmos-labs/desmos/v3/x/relationships/types"
)

type Relationship struct {
	relationshipstypes.Relationship
	Height int64
}

func NewRelationship(relationship relationshipstypes.Relationship, height int64) Relationship {
	return Relationship{
		Relationship: relationship,
		Height:       height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

type Blockage struct {
	relationshipstypes.UserBlock
	Height int64
}

func NewBlockage(blockage relationshipstypes.UserBlock, height int64) Blockage {
	return Blockage{
		UserBlock: blockage,
		Height:    height,
	}
}
