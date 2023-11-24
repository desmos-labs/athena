package notifications

import (
	"fmt"

	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
)

// BuildMessage builds the notification message that should be sent based on the given recipient, notification and data
func (m *Module) BuildMessage(recipient types.NotificationRecipient, config types.NotificationData) (types.NotificationMessage, error) {
	var androidConfig *messaging.AndroidConfig
	var apnsConfig *messaging.APNSConfig
	var webpushConfig *messaging.WebpushConfig

	if dataWithConfig, ok := config.(types.NotificationDataWithConfig); ok {
		androidConfig = dataWithConfig.GetAndroidConfig()
		apnsConfig = dataWithConfig.GetAPNSConfig()
		webpushConfig = dataWithConfig.GetWebpushConfig()
	}

	switch recipient := recipient.(type) {
	case *types.NotificationTopicRecipient:
		return types.NewSingleNotificationMessage(&messaging.Message{
			Topic:        recipient.Topic,
			Data:         config.GetAdditionalData(),
			Notification: config.GetNotification(),
			Android:      androidConfig,
			APNS:         apnsConfig,
			Webpush:      webpushConfig,
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
			Data:         config.GetAdditionalData(),
			Notification: config.GetNotification(),
			Android:      androidConfig,
			APNS:         apnsConfig,
			Webpush:      webpushConfig,
		}), nil

	default:
		return nil, fmt.Errorf("invalid notification recipient: %T", recipient)
	}
}
