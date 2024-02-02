package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/desmos-labs/athena/v2/types"
)

// GetTokens returns the tokens to be used in order to send the notification to the devices of the given recipient
func (m *Module) getUserTokens(recipient string) ([]string, error) {
	tokens, err := m.db.GetUserTokens(recipient)
	if err != nil {
		return nil, err
	}

	// Extract the tokens values
	tokensValues := make([]string, len(tokens))
	for i, token := range tokens {
		tokensValues[i] = token.Token
	}
	return tokensValues, nil
}

// SendAndStoreNotification sends the given notification to the given recipient and stores it inside the database.
func (m *Module) SendAndStoreNotification(recipient types.NotificationRecipient, notification types.NotificationData) error {
	// Send the notification
	err := m.notificationSender(recipient, notification)
	if err != nil {
		return err
	}

	// Store the notification (if enabled)
	if m.cfg.PersistHistory {
		return m.db.SaveNotification(recipient, notification)
	}

	return nil
}

// sendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
func (m *Module) sendNotification(recipient types.NotificationRecipient, notification types.NotificationData) error {
	if _, hasRecipientField := notification.GetAdditionalData()[types.RecipientKey]; !hasRecipientField {
		notification.GetAdditionalData()[types.RecipientKey] = recipient.String()
	}

	// Build the message
	message, err := m.buildMessage(recipient, notification)
	if err != nil {
		return err
	}

	// Context with 5 seconds to send the notification
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the message
	switch notificationMessage := message.(type) {
	case *types.SingleNotificationMessage:
		_, err = m.client.Send(ctx, notificationMessage.Message)
	case *types.MultiNotificationMessage:
		_, err = m.client.SendEachForMulticast(ctx, notificationMessage.MulticastMessage)
	}
	if err != nil {
		return fmt.Errorf("error while sending notification: %s", err)
	}

	return nil
}
