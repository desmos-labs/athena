package notifications

import (
	"strings"

	"github.com/rs/zerolog/log"

	poststypes "github.com/desmos-labs/desmos/v7/x/posts/types"

	"github.com/desmos-labs/athena/v2/types"
	notificationsbuilder "github.com/desmos-labs/athena/v2/x/notifications/builder"
)

func (m *Module) getPostNotificationData(originalPost types.Post, reply types.Post, builder notificationsbuilder.PostNotificationBuilder) types.NotificationData {
	if builder == nil {
		return nil
	}
	return builder(originalPost, reply)
}

// -------------------------------------------------------------------------------------------------------------------

// SendPostNotifications sends all the notifications to the users that are somehow involved with the given post.
// These include:
// - the author of the original post to which the post is a reply (if any)
// - the users mentioned inside the post
// - the authors of the various referenced posts (if this post is a reply/repost/quote)
func (m *Module) SendPostNotifications(height int64, subspaceID uint64, postID uint64) error {
	post, err := m.postsModule.GetPost(height, subspaceID, postID)
	if err != nil {
		return err
	}

	// List of users already notified
	var notifiedUsers []string

	// Send post references notifications
	for _, reference := range post.ReferencedPosts {
		// Do nothing if the post with the same id is both the original post and the post to which has been replied
		if reference.PostID == post.ConversationID {
			continue
		}

		originalPost, err := m.postsModule.GetPost(height, subspaceID, reference.PostID)
		if err != nil {
			return err
		}

		err = m.sendPostReferenceNotification(originalPost, reference.Type, post, notifiedUsers)
		if err != nil {
			return err
		}

		notifiedUsers = append(notifiedUsers, originalPost.Author)
	}

	// Send conversation notification
	if post.ConversationID != 0 {
		conversationPost, err := m.postsModule.GetPost(height, subspaceID, post.ConversationID)
		if err != nil {
			return err
		}

		err = m.sendConversationNotification(conversationPost, post, notifiedUsers)
		if err != nil {
			return err
		}

		notifiedUsers = append(notifiedUsers, conversationPost.Author)
	}

	// Send mentions notification
	if post.Entities != nil {
		for _, mention := range post.Entities.Mentions {
			err = m.sendPostMentionNotification(post, mention, notifiedUsers)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Module) getMentionNotificationData(post types.Post, mention poststypes.TextTag, builder notificationsbuilder.MentionNotificationBuilder) types.NotificationData {
	if builder == nil {
		return nil
	}
	return builder(post, mention)
}

func (m *Module) sendConversationNotification(originalPost types.Post, reply types.Post, notifiedUsers []string) error {
	// Skip if the post author and the original author are the same
	if originalPost.Author == reply.Author {
		return nil
	}

	// Skip if the referenced post author has already been notified
	if hasBeenNotified(originalPost.Author, notifiedUsers) {
		return nil
	}

	// Get the notification data
	data := m.getPostNotificationData(originalPost, reply, m.notificationsBuilder.Posts().Comment())
	if data == nil {
		return nil
	}

	log.Trace().Str("module", m.Name()).Str("recipient", originalPost.Author).
		Str("notification type", "post reply").Msg("sending notification")

	return m.SendAndStoreNotification(types.NewNotificationUserRecipient(originalPost.Author), data)
}

func (m *Module) sendPostReferenceNotification(originalPost types.Post, referenceType poststypes.PostReferenceType, reference types.Post, notifiedUsers []string) error {
	// Skip if the referenced post and the original post authors are the same
	if reference.Author == originalPost.Author {
		return nil
	}

	// Skip if the referenced post author has already been notified
	if hasBeenNotified(originalPost.Author, notifiedUsers) {
		return nil
	}

	var data types.NotificationData
	switch referenceType {
	case poststypes.POST_REFERENCE_TYPE_REPLY:
		data = m.getPostNotificationData(originalPost, reference, m.notificationsBuilder.Posts().Reply())

	case poststypes.POST_REFERENCE_TYPE_REPOST:
		data = m.getPostNotificationData(originalPost, reference, m.notificationsBuilder.Posts().Repost())

	case poststypes.POST_REFERENCE_TYPE_QUOTE:
		data = m.getPostNotificationData(originalPost, reference, m.notificationsBuilder.Posts().Quote())
	}

	if data == nil {
		return nil
	}

	log.Trace().Str("module", m.Name()).Str("recipient", originalPost.Author).
		Str("notification type", "post reference").Msg("sending notification")

	return m.SendAndStoreNotification(types.NewNotificationUserRecipient(originalPost.Author), data)
}

func (m *Module) sendPostMentionNotification(post types.Post, mention poststypes.TextTag, notifiedUsers []string) error {
	// Skip if the post author and the mentioned user is the same
	if post.Author == mention.Tag {
		return nil
	}

	// Skip if the mentioned user has already been notified
	if hasBeenNotified(mention.Tag, notifiedUsers) {
		return nil
	}

	data := m.getMentionNotificationData(post, mention, m.notificationsBuilder.Posts().Mention())
	if data == nil {
		return nil
	}

	log.Trace().Str("module", m.Name()).Str("recipient", mention.Tag).
		Str("notification type", "post mention").Msg("sending notification")

	return m.SendAndStoreNotification(types.NewNotificationUserRecipient(mention.Tag), data)
}

func hasBeenNotified(user string, notifiedUsers []string) bool {
	for _, notifiedUser := range notifiedUsers {
		if strings.EqualFold(user, notifiedUser) {
			return true
		}
	}
	return false
}
