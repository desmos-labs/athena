package types

import (
	"time"
)

type Notification struct {
	RecipientAddress string
	Type             string
	Data             map[string]string
	Timestamp        time.Time
}

func NewNotification(recipient string, notificationType string, data map[string]string, timestamp time.Time) Notification {
	return Notification{
		RecipientAddress: recipient,
		Type:             notificationType,
		Data:             data,
		Timestamp:        timestamp,
	}
}
