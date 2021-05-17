package notifications

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/desmos-labs/djuno/config"
	"github.com/desmos-labs/djuno/notifications"
)

// setupNotifications allows to properly setup the Firebase Cloud Messaging client so that
// it can later be used to send push notifications to the subscribing devices.
func setupNotifications(cfg *config.NotificationsConfig) error {
	firebaseCfg := firebase.Config{ProjectID: cfg.FirebaseProjectID}

	// Build the firebase app
	app, err := firebase.NewApp(context.Background(), &firebaseCfg, option.WithCredentialsFile(cfg.FirebaseCredentialsFile))
	if err != nil {
		return err
	}

	// Build the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	notifications.MsgClient = client
	return nil
}
