package notifications

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SaveNotification(recipient types.NotificationRecipient, notification types.NotificationData) error
	SaveToken(token types.NotificationToken) error
	GetUserTokens(userAddress string) ([]types.NotificationToken, error)
}
