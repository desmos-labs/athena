package notifications

import (
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/v2/types"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

func (m *Module) getReactionNotificationData(post types.Post, reaction types.Reaction, builder notificationsbuilder.ReactionNotificationBuilder) *notificationsbuilder.NotificationData {
	if builder == nil {
		return nil
	}
	return builder(post, reaction)
}

// -------------------------------------------------------------------------------------------------------------------

// SendReactionNotifications sends all the notifications to the author of the post that has been reacted to
func (m *Module) SendReactionNotifications(reaction types.Reaction) error {
	// Get the post
	post, err := m.postsModule.GetPost(reaction.Height, reaction.SubspaceID, reaction.PostID)
	if err != nil {
		return err
	}

	// Skip if the user reacting and the post author are the same
	if post.Author == reaction.Author {
		return nil
	}

	data := m.getReactionNotificationData(post, reaction, m.notificationBuilder.Reactions().Reaction())
	if data == nil {
		return nil
	}

	log.Debug().Str("module", m.Name()).Str("recipient", post.Author).
		Str("notification type", notificationsbuilder.TypeReaction).Msg("sending notification")
	return m.sendNotification(post.Author, data.Notification, data.Data)
}
