package devices

import (
	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"
	"github.com/desmos-labs/djuno/v2/x/notifications"
	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"
)

// --------------------------------------------------------------------------------------------------------------------

var (
	_ messagebuilder.FirebaseMessageBuilder = &MessageBuilder{}
)

// MessageBuilder represents the structure that allows building a Firebase message that is sent to specific devices
type MessageBuilder struct {
	db notifications.Database
}

func NewMessageBuilder(db notifications.Database) *MessageBuilder {
	return &MessageBuilder{
		db: db,
	}
}

// BuildMessage implements FirebaseMessageBuilder
func (b *MessageBuilder) BuildMessage(recipient string, config *messagebuilder.MessageConfig) (types.NotificationMessage, error) {
	tokens, err := b.db.GetUserTokens(recipient)
	if err != nil {
		return nil, err
	}

	// Extract the tokens values
	tokensValues := make([]string, len(tokens))
	for i, token := range tokens {
		tokensValues[i] = token.Token
	}

	if len(tokensValues) == 0 {
		return nil, nil
	}

	return types.NewMultiNotificationMessage(&messaging.MulticastMessage{
		Tokens:       tokensValues,
		Data:         config.Data,
		Notification: config.Notification,
		Android:      config.Android,
		Webpush:      config.WebPush,
		APNS:         config.Apple,
	}), nil
}
