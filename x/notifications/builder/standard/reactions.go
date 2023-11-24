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
	m UtilityModule
}

func NewDefaultReactionsNotificationsBuilder(utilityModule UtilityModule) *DefaultReactionsNotificationsBuilder {
	return &DefaultReactionsNotificationsBuilder{
		m: utilityModule,
	}
}

func (d DefaultReactionsNotificationsBuilder) Reaction() notificationsbuilder.ReactionNotificationBuilder {
	return func(post types.Post, reaction types.Reaction) types.NotificationData {
		return types.NewStdNotificationDataWithConfig(
			&messaging.Notification{
				Title: "Someone reacted to your post! 🎉",
				Body:  fmt.Sprintf("%s reacted to your post", d.m.GetDisplayName(reaction.Author)),
			},
			map[string]string{
				types.NotificationTypeKey:   types.TypeReaction,
				types.NotificationActionKey: types.ActionOpenPost,

				types.SubspaceIDKey:     fmt.Sprintf("%d", post.SubspaceID),
				types.PostIDKey:         fmt.Sprintf("%d", post.ID),
				types.ReactionIDKey:     fmt.Sprintf("%d", reaction.ID),
				types.ReactionAuthorKey: reaction.Author,
			},
		)
	}
}
