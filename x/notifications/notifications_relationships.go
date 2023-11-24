package notifications

import (
	"github.com/desmos-labs/djuno/v2/types"
	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"

	"github.com/rs/zerolog/log"
)

func (m *Module) getRelationshipNotificationData(relationship types.Relationship, builder notificationsbuilder.RelationshipNotificationBuilder) types.NotificationData {
	if builder == nil {
		return nil
	}
	return builder(relationship)
}

// -------------------------------------------------------------------------------------------------------------------

// SendRelationshipNotifications sends the notification to the user towards which a relationship has just been created
func (m *Module) SendRelationshipNotifications(relationship types.Relationship) error {
	// Skip if the user and the counterparty are the same
	if relationship.Creator == relationship.Counterparty {
		return nil
	}

	data := m.getRelationshipNotificationData(relationship, m.notificationsBuilder.Relationships().Relationship())
	if data == nil {
		return nil
	}

	log.Trace().Str("module", m.Name()).Str("recipient", relationship.Counterparty).
		Str("notification type", "relationship").Msg("sending notification")

	return m.SendAndStoreNotification(types.NewNotificationUserRecipient(relationship.Counterparty), data)
}
