package notifications

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SaveNotification(notification types.Notification) error
	SaveToken(token types.NotificationToken) error
	GetUserTokens(userAddress string) ([]types.NotificationToken, error)
}
