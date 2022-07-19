package notifications

import (
	"fmt"

	"firebase.google.com/go/messaging"
)

// SendRelationshipNotifications sends the notification to the user towards which a relationship has just been created
func (m *Module) SendRelationshipNotifications(subspaceID uint64, user, counterparty string) error {
	// Skip if the user and the counterparty are the same
	if user == counterparty {
		return nil
	}

	notification := &messaging.Notification{
		Title: "You have a new follower! 👥",
		Body:  fmt.Sprintf("%s has started following you", user),
	}

	data := map[string]string{
		NotificationTypeKey:   TypeFollow,
		NotificationActionKey: ActionOpenProfile,

		SubspaceIDKey:          fmt.Sprintf("%d", subspaceID),
		RelationshipCreatorKey: user,
	}

	return m.sendNotification(counterparty, notification, data)
}
