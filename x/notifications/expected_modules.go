package notifications

import (
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/athena/v2/types"
)

type ProfilesModule interface {
	GetUserProfile(userAddress string) (*types.Profile, error)
}

type PostsModule interface {
	GetPost(height int64, subspaceID uint64, postID uint64) (types.Post, error)
}

type ReactionsModule interface {
	GetReactionID(tx *juno.Tx, index int) (uint32, error)
	GetReaction(height int64, subspaceID uint64, postID uint64, reactionID uint32) (types.Reaction, error)
}
