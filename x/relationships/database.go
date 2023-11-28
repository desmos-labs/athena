package relationships

import (
	"github.com/desmos-labs/athena/types"
)

type Database interface {
	SaveRelationship(relationship types.Relationship) error
	DeleteRelationship(relationship types.Relationship) error
	DeleteAllRelationships(height int64, subspaceID uint64) error
	SaveUserBlock(block types.Blockage) error
	DeleteBlockage(block types.Blockage) error
	DeleteAllUserBlocks(height int64, subspaceID uint64) error
}
