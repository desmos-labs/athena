package types

import (
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
)

type Post struct {
	poststypes.Post
	Height int64
}

func NewPost(post poststypes.Post, height int64) Post {
	return Post{
		Post:   post,
		Height: height,
	}
}

type PostAttachment struct {
	poststypes.Attachment
	Height int64
}

func NewPostAttachment(attachment poststypes.Attachment, height int64) PostAttachment {
	return PostAttachment{
		Attachment: attachment,
		Height:     height,
	}
}

type PollAnswer struct {
	poststypes.UserAnswer
	Height int64
}

func NewPollAnswer(answer poststypes.UserAnswer, height int64) PollAnswer {
	return PollAnswer{
		UserAnswer: answer,
		Height:     height,
	}
}

type PostsParams struct {
	poststypes.Params
	Height int64
}

func NewPostsParams(params poststypes.Params, height int64) PostsParams {
	return PostsParams{
		Params: params,
		Height: height,
	}
}
