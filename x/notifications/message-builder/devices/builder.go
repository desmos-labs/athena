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

// GetTokens returns the tokens to be used in order to send the notification to the devices of the given recipient
func (b *MessageBuilder) GetTokens(recipient string) ([]string, error) {
	tokens, err := b.db.GetUserTokens(recipient)
	if err != nil {
		return nil, err
	}

	// Extract the tokens values
	tokensValues := make([]string, len(tokens))
	for i, token := range tokens {
		tokensValues[i] = token.Token
	}
	return tokensValues, nil
}

// BuildMessage implements FirebaseMessageBuilder
func (b *MessageBuilder) BuildMessage(recipient string, config *messagebuilder.MessageConfig) (types.NotificationMessage, error) {
	tokens, err := b.GetTokens(recipient)
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, nil
	}

	return types.NewMultiNotificationMessage(&messaging.MulticastMessage{
		Tokens:       tokens,
		Data:         config.Data,
		Notification: config.Notification,
	}), nil
}
