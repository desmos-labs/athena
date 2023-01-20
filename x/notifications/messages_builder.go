package notifications

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
)

// BuildMessage builds the notification message that should be sent based on the given recipient, notification and data
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
