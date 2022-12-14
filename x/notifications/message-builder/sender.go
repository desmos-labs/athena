package messagebuilder

import (
	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"

	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
)

// FirebaseMessageBuilderCreator represents a FirebaseMessageBuilder creator
type FirebaseMessageBuilderCreator func(ctx notificationscontext.Context) FirebaseMessageBuilder

type MessageConfig struct {
	Data         map[string]string
	Notification *messaging.Notification
	Android      *messaging.AndroidConfig
	WebPush      *messaging.WebpushConfig
	Apple        *messaging.APNSConfig
}

// FirebaseMessageBuilder represents the interface that allows to build a Firebase message
type FirebaseMessageBuilder interface {
	// BuildMessage builds the Firebase message to be sent, based on the given data
	BuildMessage(recipient string, config *MessageConfig) (types.NotificationMessage, error)
}
