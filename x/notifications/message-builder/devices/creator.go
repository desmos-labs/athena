package devices

import (
	"fmt"

	"github.com/desmos-labs/djuno/v2/x/notifications"
	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
	messagebuilder "github.com/desmos-labs/djuno/v2/x/notifications/message-builder"
)

// CreateMessageBuilder creates a new FirebaseMessageBuilder instance that sends out topic messages
func CreateMessageBuilder(ctx notificationscontext.Context) messagebuilder.FirebaseMessageBuilder {
	db, ok := ctx.Database.(notifications.Database)
	if !ok {
		panic(fmt.Errorf("invalid database type: expected notifications.Database, got %T instead", ctx.Database))
	}
	return NewMessageBuilder(db)
}
