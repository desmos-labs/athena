package posts

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SavePost(post types.Post) error
	DeletePost(height int64, subspaceID uint64, postID uint64) error
	DeleteAllPosts(height int64, subspaceID uint64) error
	SavePostTx(tx types.PostTransaction) error
	SavePostAttachment(attachment types.PostAttachment) error
	DeletePostAttachment(height int64, subspaceID uint64, postID uint64, attachmentID uint32) error
	SavePollAnswer(answer types.PollAnswer) error
	SavePostsParams(params types.PostsParams) error
}
