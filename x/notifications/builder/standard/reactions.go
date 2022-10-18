package standard

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

var (
	_ notificationsbuilder.ReactionsNotificationsBuilder = &DefaultReactionsNotificationsBuilder{}
)

type DefaultReactionsNotificationsBuilder struct {
	m notificationsbuilder.UtilityModule
}

func NewDefaultReactionsNotificationsBuilder(utilityModule notificationsbuilder.UtilityModule) *DefaultReactionsNotificationsBuilder {
	return &DefaultReactionsNotificationsBuilder{
		m: utilityModule,
	}
}

func (d DefaultReactionsNotificationsBuilder) Reaction() notificationsbuilder.ReactionNotificationBuilder {
	return func(post types.Post, reaction types.Reaction) *notificationsbuilder.NotificationData {
		return &notificationsbuilder.NotificationData{
			Notification: &messaging.Notification{
				Title: "Someone reacted to your post! 🎉",
				Body:  fmt.Sprintf("%s reacted to your post", d.m.GetDisplayName(reaction.Author)),
			},
			Data: map[string]string{
				notificationsbuilder.NotificationTypeKey:   notificationsbuilder.TypeReaction,
				notificationsbuilder.NotificationActionKey: notificationsbuilder.ActionOpenPost,

				notificationsbuilder.SubspaceIDKey:     fmt.Sprintf("%d", post.SubspaceID),
				notificationsbuilder.PostIDKey:         fmt.Sprintf("%d", post.ID),
				notificationsbuilder.ReactionAuthorKey: reaction.Author,
			},
		}
	}
}
