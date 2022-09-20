package notifications

import (
	"fmt"
	"strings"

	"firebase.google.com/go/v4/messaging"
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

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

func (m *Module) sendConversationNotification(originalPost types.Post, reply types.Post, notifiedUsers []string) error {
	// Skip if the post author and the original author are the same
	if originalPost.Author == reply.Author {
		return nil
	}

	// Skip if the referenced post author has already been notified
	if hasBeenNotified(originalPost.Author, notifiedUsers) {
		return nil
	}

	notification := &messaging.Notification{
		Title: "Someone replied to your post! 💬",
		Body:  fmt.Sprintf("%s replied to your post", reply.Author),
	}

	data := map[string]string{
		NotificationTypeKey:   TypeReply,
		NotificationActionKey: ActionOpenPost,

		SubspaceIDKey: fmt.Sprintf("%d", reply.SubspaceID),
		PostIDKey:     fmt.Sprintf("%d", reply.ID),
		PostAuthorKey: reply.Author,
	}

	return m.sendNotification(originalPost.Author, notification, data)
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

	var notificationType string
	var notification *messaging.Notification
	switch referenceType {
	case poststypes.POST_REFERENCE_TYPE_REPLY:
		notificationType = TypeReply
		notification = &messaging.Notification{
			Title: "Someone commented your post! 💬",
			Body:  fmt.Sprintf("%s commented on your post", reference.Author),
		}

	case poststypes.POST_REFERENCE_TYPE_REPOST:
		notificationType = TypeRepost
		notification = &messaging.Notification{
			Title: "Someone reposted your post! 💬",
			Body:  fmt.Sprintf("%s reposted your post", reference.Author),
		}

	case poststypes.POST_REFERENCE_TYPE_QUOTE:
		notificationType = TypeQuote
		notification = &messaging.Notification{
			Title: "Someone quoted your post! 💬",
			Body:  fmt.Sprintf("%s quoted your post", reference.Author),
		}
	}

	data := map[string]string{
		NotificationTypeKey:   notificationType,
		NotificationActionKey: ActionOpenPost,

		SubspaceIDKey: fmt.Sprintf("%d", originalPost.SubspaceID),
		PostIDKey:     fmt.Sprintf("%d", originalPost.ID),
		PostAuthorKey: reference.Author,
	}

	return m.sendNotification(originalPost.Author, notification, data)
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

	notification := &messaging.Notification{
		Title: "Someone mentioned you inside a post! 💬",
		Body:  fmt.Sprintf("%s mentioned you post", post.Author),
	}

	data := map[string]string{
		NotificationTypeKey:   TypeMention,
		NotificationActionKey: ActionOpenPost,

		SubspaceIDKey: fmt.Sprintf("%d", post.SubspaceID),
		PostIDKey:     fmt.Sprintf("%d", post.ID),
		PostAuthorKey: post.Author,
	}

	return m.sendNotification(mention.Tag, notification, data)
}

func hasBeenNotified(user string, notifiedUsers []string) bool {
	for _, notifiedUser := range notifiedUsers {
		if strings.EqualFold(user, notifiedUser) {
			return true
		}
	}
	return false
}
