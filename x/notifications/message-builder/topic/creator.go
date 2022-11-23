package topic

import (
	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"
)

// CreateMessageBuilder creates a new FirebaseMessageBuilder instance that sends out topic messages
func CreateMessageBuilder(_ notificationscontext.Context) messagebuilder.FirebaseMessageBuilder {
	return NewMessageBuilder()
}
