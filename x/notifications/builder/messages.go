package builder

import (
	"github.com/desmos-labs/athena/v2/types"
)

// MessagesBuilder represents a NotificationMessage builder
type MessagesBuilder = func(recipient types.NotificationRecipient, data types.NotificationData) (types.NotificationMessage, error)
