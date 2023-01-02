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
	m UtilityModule
}

func NewDefaultRelationshipsNotificationsBuilder(utilityModule UtilityModule) *DefaultRelationshipsNotificationsBuilder {
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
				types.NotificationTypeKey:   types.TypeFollow,
				types.NotificationActionKey: types.ActionOpenProfile,

				types.SubspaceIDKey:          fmt.Sprintf("%d", relationship.SubspaceID),
				types.RelationshipCreatorKey: relationship.Creator,
			},
		}
	}
}
