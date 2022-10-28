package types

import (
	"time"
)

type Notification struct {
	RecipientAddress string
	Data             map[string]string
	Timestamp        time.Time
}

func NewNotification(recipient string, data map[string]string, timestamp time.Time) Notification {
	return Notification{
		RecipientAddress: recipient,
		Data:             data,
		Timestamp:        timestamp,
	}
}
