package notifications

import (
	"github.com/desmos-labs/djuno/v2/types"
)

type Database interface {
	SaveNotification(notification types.Notification) error
}
