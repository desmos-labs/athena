package posts

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v3/node/remote"

	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshPostsData refreshes all the posts' data for the given subspace
func (m *Module) RefreshPostsData(height int64, subspaceID uint64) error {
	posts, err := m.QuerySubspacePosts(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllPosts(height)
	if err != nil {
		return fmt.Errorf("error while deleting all the posts: %s", err)
	}

	// Refresh posts
	for _, post := range posts {
		log.Info().Uint64("subspace", post.SubspaceID).Uint64("post", post.ID).Msg("refreshing post")

		err = m.db.SavePost(post)
		if err != nil {
			return fmt.Errorf("error while saving post: %s", err)
		}

		attachments, err := m.queryPostAttachments(height, post.SubspaceID, post.ID)
		if err != nil {
			return err
		}

		// Refresh attachments
		for _, attachment := range attachments {
			log.Info().Uint64("subspace", attachment.SubspaceID).Uint64("post", attachment.PostID).
				Uint32("attachment", attachment.ID).Msg("refreshing attachment")

			err = m.db.SavePostAttachment(attachment)
			if err != nil {
				return fmt.Errorf("error while saving post attachment: %s", err)
			}

			// Refresh poll answers
			if _, isPoll := attachment.Content.GetCachedValue().(*poststypes.Poll); isPoll {
				log.Info().Uint64("subspace", attachment.SubspaceID).Uint64("post", attachment.PostID).
					Uint32("poll", attachment.ID).Msg("refreshing poll answers")

				answers, err := m.queryPollAnswers(height, attachment.SubspaceID, attachment.PostID, attachment.ID)
				if err != nil {
					return err
				}

				for _, answer := range answers {
					err = m.db.SavePollAnswer(answer)
					if err != nil {
						return fmt.Errorf("errror while saving poll answer: %s", err)
					}
				}
			}
		}
	}

	return nil
}

// QuerySubspacePosts queries all the posts present inside the given subspace at the provided height
func (m *Module) QuerySubspacePosts(height int64, subspaceID uint64) ([]types.Post, error) {
	var posts []types.Post

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.SubspacePosts(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QuerySubspacePostsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, post := range res.Posts {
			posts = append(posts, types.NewPost(post, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return posts, nil
}

// QuerySubspacePosts queries all the attachments for the given post at the provided height
func (m *Module) queryPostAttachments(height int64, subspaceID uint64, postID uint64) ([]types.PostAttachment, error) {
	var attachments []types.PostAttachment

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.PostAttachments(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QueryPostAttachmentsRequest{
				SubspaceId: subspaceID,
				PostId:     postID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, attachment := range res.Attachments {
			attachments = append(attachments, types.NewPostAttachment(attachment, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return attachments, nil
}

// QuerySubspacePosts queries all the posts present inside the given subspace at the provided height
func (m *Module) queryPollAnswers(height int64, subspaceID uint64, postID uint64, pollID uint32) ([]types.PollAnswer, error) {
	var answers []types.PollAnswer

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.PollAnswers(
			remote.GetHeightRequestContext(context.Background(), height),
			&poststypes.QueryPollAnswersRequest{
				SubspaceId: subspaceID,
				PostId:     postID,
				PollId:     pollID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, answer := range res.Answers {
			answers = append(answers, types.NewPollAnswer(answer, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return answers, nil
}
