package notifications

import (
	"context"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

var (
	msgClient *messaging.Client
)

const (
	ProjectID = "mooncake-16aaa"

	AndroidChannelID = "mooncake_posts"

	ClickActionKey   = "click_action"
	ClickActionValue = "FLUTTER_NOTIFICATION_CLICK"
)

// SetupFirebase allows to properly setup the Firebase Cloud Messaging client so that
// it can later be used to send push notifications to the subscribing devices.
func SetupFirebase(credentialsFile string) error {
	config := firebase.Config{ProjectID: ProjectID}

	// Build the firebase app
	app, err := firebase.NewApp(context.Background(), &config, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return err
	}

	// Build the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	msgClient = client
	return nil
}

// SendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
// If some error rises during the process, it is returned.
func SendNotification(topic string, notification *messaging.Notification, data map[string]string) error {
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
	_, err := msgClient.Send(ctx, &message)
	return err
}
