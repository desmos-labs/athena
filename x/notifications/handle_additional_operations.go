package notifications

import (
	"context"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m Module) RunAdditionalOperations() error {
	firebaseCfg := firebase.Config{ProjectID: m.cfg.FirebaseProjectID}

	// Build the firebase app
	app, err := firebase.NewApp(
		context.Background(),
		&firebaseCfg,
		option.WithCredentialsFile(m.cfg.FirebaseCredentialsFile),
	)
	if err != nil {
		return err
	}

	// Build the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	MsgClient = client
	return nil
}
