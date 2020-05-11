package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// PostReaction represents a reaction that is added to a post
type PostReaction struct {
	PostID    posts.PostID
	Value     string
	ShortCode string
	User      sdk.AccAddress
}

// NewPostReaction returns a new PostReaction object containing the specified values
func NewPostReaction(postID posts.PostID, value, shortCode string, user sdk.AccAddress) PostReaction {
	return PostReaction{
		PostID:    postID,
		Value:     value,
		ShortCode: shortCode,
		User:      user,
	}
}
