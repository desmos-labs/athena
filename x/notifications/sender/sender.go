package sender

import (
	"github.com/desmos-labs/djuno/v2/types"
	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
)

// NotificationSender represents a function that allows to send a notification to the given recipient
type NotificationSender = func(recipient types.NotificationRecipient, notification types.NotificationData) error

// NotificationsSenderCreator represents a function that allows to create a new NotificationSender
type NotificationsSenderCreator = func(context notificationscontext.Context) NotificationSender
