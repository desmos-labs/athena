package reactions

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	HasPost(height int64, subspaceID uint64, postID uint64) (bool, error)

	SaveReaction(reaction types.Reaction) error
	DeleteReaction(height int64, subspaceID uint64, postID uint64, reactionID uint32) error
	DeleteAllReactions(height int64, subspaceID uint64, postID uint64) error
	SaveRegisteredReaction(reaction types.RegisteredReaction) error
	DeleteRegisteredReaction(height int64, subspaceID uint64, reactionID uint32) error
	DeleteAllRegisteredReactions(height int64, subspaceID uint64) error
	SaveReactionParams(params types.ReactionParams) error
}
