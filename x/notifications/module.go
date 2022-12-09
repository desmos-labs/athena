package notifications

import (
	"context"
	"fmt"
	"time"

	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/cosmos/cosmos-sdk/codec"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/types/config"
	"google.golang.org/api/option"

	notificationsbuilder "github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

var (
	_ modules.Module            = &Module{}
	_ modules.TransactionModule = &Module{}
	_ modules.MessageModule     = &Module{}
)

type Module struct {
	cdc codec.Codec
	db  Database

	cfg    *Config
	app    *firebase.App
	client *messaging.Client

	postsModule     PostsModule
	reactionsModule ReactionsModule

	notificationBuilder notificationsbuilder.NotificationsBuilder
	messageBuilder      messagebuilder.FirebaseMessageBuilder
}

// NewModule returns a new Module instance
func NewModule(
	junoCfg config.Config,
	postsModule PostsModule, reactionsModule ReactionsModule,
	cdc codec.Codec, db Database,
) *Module {
	bz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}

	cfg, err := ParseConfig(bz)
	if err != nil {
		panic(err)
	}

	if cfg == nil {
		return nil
	}

	firebaseCfg := firebase.Config{ProjectID: cfg.FirebaseProjectID}

	// Build the firebase app
	app, err := firebase.NewApp(context.Background(), &firebaseCfg, option.WithCredentialsFile(cfg.FirebaseCredentialsFilePath))
	if err != nil {
		panic(err)
	}

	// Build the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		panic(err)
	}

	return &Module{
		cdc:             cdc,
		db:              db,
		cfg:             cfg,
		app:             app,
		client:          client,
		postsModule:     postsModule,
		reactionsModule: reactionsModule,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "notifications"
}

// WithNotificationsBuilder sets the given builder as the notifications builder
func (m *Module) WithNotificationsBuilder(builder notificationsbuilder.NotificationsBuilder) *Module {
	if builder != nil {
		m.notificationBuilder = builder
	}
	return m
}

// WithFirebaseMessageBuilder sets the given builder as the Firebase message builder
func (m *Module) WithFirebaseMessageBuilder(builder messagebuilder.FirebaseMessageBuilder) *Module {
	if builder != nil {
		m.messageBuilder = builder
	}
	return m
}

// SendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
func (m *Module) SendNotification(recipient string, notification *messaging.Notification, data map[string]string) error {
	// Set the default Flutter click action
	data[notificationsbuilder.ClickActionKey] = notificationsbuilder.ClickActionValue
	data[notificationsbuilder.RecipientKey] = recipient

	// Build the Android config
	var androidConfig *messaging.AndroidConfig
	if notification != nil {
		androidConfig = &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{ChannelID: m.cfg.AndroidChannelID},
		}
	}

	// Build the message
	message, err := m.messageBuilder.BuildMessage(recipient, &messagebuilder.MessageConfig{
		Data:         data,
		Notification: notification,
		Android:      androidConfig,
	})
	if err != nil {
		return fmt.Errorf("error while building notification message: %s", err)
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
			data[notificationsbuilder.NotificationTypeKey],
			data,
			time.Now(),
		))
	}

	return nil
}
