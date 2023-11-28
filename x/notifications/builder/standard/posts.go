package standard

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"

	"github.com/desmos-labs/athena/types"
	notificationsbuilder "github.com/desmos-labs/athena/x/notifications/builder"
)

var (
	_ notificationsbuilder.PostsNotificationsBuilder = &DefaultPostsNotificationsBuilder{}
)

type DefaultPostsNotificationsBuilder struct {
	m UtilityModule
}

func NewDefaultPostsNotificationsBuilder(utilityModule UtilityModule) *DefaultPostsNotificationsBuilder {
	return &DefaultPostsNotificationsBuilder{
		m: utilityModule,
	}
}

func (d DefaultPostsNotificationsBuilder) Comment() notificationsbuilder.PostNotificationBuilder {
	return func(originalPost types.Post, comment types.Post) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone commented your post! 💬",
				Body:  fmt.Sprintf("%s commented on your post", d.m.GetDisplayName(comment.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeComment,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey:    fmt.Sprintf("%d", originalPost.SubspaceID),
				types.PostIDKey:        fmt.Sprintf("%d", originalPost.ID),
				types.CommentIDKey:     fmt.Sprintf("%d", comment.ID),
				types.CommentAuthorKey: comment.Author,
			},
		)
	}
}

func (d DefaultPostsNotificationsBuilder) Reply() notificationsbuilder.PostNotificationBuilder {
	return func(originalPost types.Post, reply types.Post) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone replied to your post! 💬",
				Body:  fmt.Sprintf("%s replied to your post", d.m.GetDisplayName(reply.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeReply,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey:  fmt.Sprintf("%d", originalPost.SubspaceID),
				types.PostIDKey:      fmt.Sprintf("%d", originalPost.ID),
				types.ReplyIDKey:     fmt.Sprintf("%d", reply.ID),
				types.ReplyAuthorKey: reply.Author,
			},
		)
	}
}

func (d DefaultPostsNotificationsBuilder) Repost() notificationsbuilder.PostNotificationBuilder {
	return func(originalPost types.Post, repost types.Post) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone reposted your post! 💬",
				Body:  fmt.Sprintf("%s reposted your post", d.m.GetDisplayName(repost.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeRepost,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey:   fmt.Sprintf("%d", originalPost.SubspaceID),
				types.PostIDKey:       fmt.Sprintf("%d", originalPost.ID),
				types.RepostIDKey:     fmt.Sprintf("%d", repost.ID),
				types.RepostAuthorKey: repost.Author,
			},
		)
	}
}

func (d DefaultPostsNotificationsBuilder) Quote() notificationsbuilder.PostNotificationBuilder {
	return func(originalPost types.Post, quote types.Post) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone quoted your post! 💬",
				Body:  fmt.Sprintf("%s quoted your post", d.m.GetDisplayName(quote.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeQuote,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey:  fmt.Sprintf("%d", originalPost.SubspaceID),
				types.PostIDKey:      fmt.Sprintf("%d", originalPost.ID),
				types.QuoteIDKey:     fmt.Sprintf("%d", quote.ID),
				types.QuoteAuthorKey: quote.Author,
			},
		)
	}
}

func (d DefaultPostsNotificationsBuilder) Mention() notificationsbuilder.MentionNotificationBuilder {
	return func(post types.Post, mention poststypes.TextTag) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone mentioned you inside a post! 💬",
				Body:  fmt.Sprintf("%s mentioned you post", d.m.GetDisplayName(post.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeMention,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey: fmt.Sprintf("%d", post.SubspaceID),
				types.PostIDKey:     fmt.Sprintf("%d", post.ID),
				types.PostAuthorKey: post.Author,
			},
		)
	}
}
