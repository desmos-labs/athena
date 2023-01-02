package notifications

import (
	"context"
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
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

// BuildMessage builds the notification message that should be sent based on the given recipient, notifiation and data
func (m *Module) BuildMessage(recipient types.NotificationRecipient, config *types.NotificationConfig) (types.NotificationMessage, error) {
	switch recipient := recipient.(type) {
	case *types.NotificationTopicRecipient:
		return types.NewSingleNotificationMessage(&messaging.Message{
			Topic:        recipient.Topic,
			Data:         config.Data,
			Notification: config.Notification,
			Android:      config.Android,
			APNS:         config.APNS,
		}), nil

	case *types.NotificationUserRecipient:
		tokens, err := m.getUserTokens(recipient.Address)
		if err != nil {
			return nil, err
		}
		if len(tokens) == 0 {
			return nil, nil
		}
		return types.NewMultiNotificationMessage(&messaging.MulticastMessage{
			Tokens:       tokens,
			Data:         config.Data,
			Notification: config.Notification,
			Android:      config.Android,
			APNS:         config.APNS,
		}), nil

	default:
		return nil, fmt.Errorf("invalid notification recipient: %T", recipient)
	}
}

// SendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
func (m *Module) SendNotification(recipient types.NotificationRecipient, config *types.NotificationConfig) error {
	if _, hasRecipientField := config.Data[types.RecipientKey]; !hasRecipientField {
		config.Data[types.RecipientKey] = recipient.String()
	}

	// Build the message
	message, err := m.BuildMessage(recipient, config)
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
		_, err = m.client.SendMulticast(ctx, notificationMessage.MulticastMessage)
	}
	if err != nil {
		return fmt.Errorf("error while sending notification: %s", err)
	}

	// Store the notification (if enabled)
	if m.cfg.PersistHistory {
		return m.db.SaveNotification(types.NewNotification(
			recipient,
			config.Type,
			config.Data,
			time.Now(),
		))
	}

	return nil
}
