package notifications

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/types/config"
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
	messagesBuilder     notificationsbuilder.MessagesBuilder
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

	// Build the module
	module := &Module{
		cdc:             cdc,
		db:              db,
		cfg:             cfg,
		app:             app,
		client:          client,
		postsModule:     postsModule,
		reactionsModule: reactionsModule,
	}

	// Set the default messages builder
	module = module.WithMessagesBuilder(module.BuildMessage)

	return module
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

// WithMessagesBuilder sets the given builder as the messages builder
func (m *Module) WithMessagesBuilder(builder notificationsbuilder.MessagesBuilder) *Module {
	if builder != nil {
		m.messagesBuilder = builder
	}
	return m
}
