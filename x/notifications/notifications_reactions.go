package notifications

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"firebase.google.com/go/v4/messaging"
)

// SendReactionNotifications sends all the notifications to the author of the post that has been reacted to
func (m *Module) SendReactionNotifications(height int64, subspaceID uint64, postID uint64, user string) error {
	// Get the post
	post, err := m.postsModule.GetPost(height, subspaceID, postID)
	if err != nil {
		return err
	}

	// Skip if the user reacting and the post author are the same
	if post.Author == user {
		return nil
	}

	notification := &messaging.Notification{
		Title: "Someone reacted to your post! 🎉",
		Body:  fmt.Sprintf("%s reacted to your post", m.getDisplayName(user)),
	}

	data := map[string]string{
		NotificationTypeKey:   TypeReaction,
		NotificationActionKey: ActionOpenPost,

		SubspaceIDKey:     fmt.Sprintf("%d", post.SubspaceID),
		PostIDKey:         fmt.Sprintf("%d", post.ID),
		ReactionAuthorKey: user,
	}

	log.Debug().Str("module", m.Name()).Str("recipient", post.Author).
		Str("notification type", TypeReaction).Msg("sending notification")
	return m.sendNotification(post.Author, notification, data)
}
