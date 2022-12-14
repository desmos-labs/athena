package types

import (
	"time"

	"firebase.google.com/go/v4/messaging"
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

// --------------------------------------------------------------------------------------------------------------------

type NotificationToken struct {
	UserAddress string
	Token       string
	Timestamp   time.Time
}

func NewNotificationToken(userAddress string, token string, timestamp time.Time) NotificationToken {
	return NotificationToken{
		UserAddress: userAddress,
		Token:       token,
		Timestamp:   timestamp,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type NotificationMessage interface {
	isNotificationMessage()
}

type SingleNotificationMessage struct {
	*messaging.Message
}

func NewSingleNotificationMessage(message *messaging.Message) *SingleNotificationMessage {
	return &SingleNotificationMessage{
		Message: message,
	}
}

func (s *SingleNotificationMessage) isNotificationMessage() {}

type MultiNotificationMessage struct {
	*messaging.MulticastMessage
}

func NewMultiNotificationMessage(message *messaging.MulticastMessage) *MultiNotificationMessage {
	return &MultiNotificationMessage{
		MulticastMessage: message,
	}
}

func (m *MultiNotificationMessage) isNotificationMessage() {}
