package standard

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

var (
	_ notificationsbuilder.RelationshipsNotificationsBuilder = &DefaultRelationshipsNotificationsBuilder{}
)

type DefaultRelationshipsNotificationsBuilder struct {
	m notificationsbuilder.UtilityModule
}

func NewDefaultRelationshipsNotificationsBuilder(utilityModule notificationsbuilder.UtilityModule) *DefaultRelationshipsNotificationsBuilder {
	return &DefaultRelationshipsNotificationsBuilder{
		m: utilityModule,
	}
}

func (d DefaultRelationshipsNotificationsBuilder) Relationship() notificationsbuilder.RelationshipNotificationBuilder {
	return func(relationship types.Relationship) *notificationsbuilder.NotificationData {
		return &notificationsbuilder.NotificationData{
			Notification: &messaging.Notification{
				Title: "You have a new follower! 👥",
				Body:  fmt.Sprintf("%s has started following you", d.m.GetDisplayName(relationship.Creator)),
			},
			Data: map[string]string{
				notificationsbuilder.NotificationTypeKey:   notificationsbuilder.TypeFollow,
				notificationsbuilder.NotificationActionKey: notificationsbuilder.ActionOpenProfile,

				notificationsbuilder.SubspaceIDKey:          fmt.Sprintf("%d", relationship.SubspaceID),
				notificationsbuilder.RelationshipCreatorKey: relationship.Creator,
			},
		}
	}
}
