package notifications

import (
	"context"
	"time"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/desmos-labs/djuno/v2/database"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/types/config"
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
	db  *database.Db

	cfg    *Config
	app    *firebase.App
	client *messaging.Client

	postsModule     PostsModule
	reactionsModule ReactionsModule

	builder notificationsbuilder.NotificationsBuilder
}

// NewModule returns a new Module instance
func NewModule(junoCfg config.Config, postsModule PostsModule, reactionsModule ReactionsModule, cdc codec.Codec, db *database.Db) *Module {
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
	m.builder = builder
	return m
}

// sendNotification allows to send to the devices subscribing to the specific topic a message
// containing the given notification and data.
func (m *Module) sendNotification(recipient string, notification *messaging.Notification, data map[string]string) error {
	// Set the default Flutter click action
	data[notificationsbuilder.ClickActionKey] = notificationsbuilder.ClickActionValue

	// Build the Android config
	var androidConfig *messaging.AndroidConfig
	if notification != nil {
		androidConfig = &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{ChannelID: m.cfg.AndroidChannelID},
		}
	}

	// Build the message
	message := messaging.Message{
		Data:         data,
		Notification: notification,
		Android:      androidConfig,
		Topic:        recipient,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send the message
	_, err := m.client.Send(ctx, &message)
	if err != nil {
		return err
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
