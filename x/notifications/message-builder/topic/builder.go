package topic

import (
	"firebase.google.com/go/v4/messaging"

	"github.com/desmos-labs/djuno/v2/types"

	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"
)

var (
	_ messagebuilder.FirebaseMessageBuilder = &MessageBuilder{}
)

// MessageBuilder represents a FirebaseMessageBuilder that builds Firebase messages that are sent to entire topics
type MessageBuilder struct {
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{}
}

// BuildMessage implements FirebaseMessageBuilder
func (m *MessageBuilder) BuildMessage(recipient string, config *messagebuilder.MessageConfig) (types.NotificationMessage, error) {
	return types.NewSingleNotificationMessage(&messaging.Message{
		Topic:        recipient,
		Data:         config.Data,
		Notification: config.Notification,
	}), nil
}
