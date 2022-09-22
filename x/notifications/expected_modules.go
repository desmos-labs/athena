package notifications

import "github.com/desmos-labs/djuno/v2/types"

type ProfilesModule interface {
	GetUserProfile(userAddress string) (*types.Profile, error)
}

type PostsModule interface {
	GetPost(height int64, subspaceID uint64, postID uint64) (types.Post, error)
}
