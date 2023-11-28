package builder

import (
	"github.com/desmos-labs/athena/types"
)

// MessagesBuilder represents a NotificationMessage builder
type MessagesBuilder = func(recipient types.NotificationRecipient, data types.NotificationData) (types.NotificationMessage, error)
