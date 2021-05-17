package utils

import (
	"context"
	"time"

	"firebase.google.com/go/messaging"
)

var (
	MsgClient *messaging.Client
)

const (
	AndroidChannelID = "mooncake_posts"

	ClickActionKey   = "click_action"
	ClickActionValue = "FLUTTER_NOTIFICATION_CLICK"
)

// SendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
// If some error rises during the process, it is returned.
func SendNotification(topic string, notification *messaging.Notification, data map[string]string) error {
	// TODO: Re-implement this
	// If disabled, just return
	//if !viper.GetBool(flags.FlagEnableNotifications) {
	//	return nil
	//}

	// Set the default Flutter click action
	data[ClickActionKey] = ClickActionValue

	// Build the Android config
	var androidConfig *messaging.AndroidConfig
	if notification != nil {
		androidConfig = &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{ChannelID: AndroidChannelID},
		}
	}

	// Build the message
	message := messaging.Message{
		Data:         data,
		Notification: notification,
		Android:      androidConfig,
		Topic:        topic,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the message
	_, err := MsgClient.Send(ctx, &message)
	return err
}
