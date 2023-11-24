package types

import (
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"
)

// NotificationRecipient represents a generic Firebase message recipient
type NotificationRecipient interface {
	GetValue() string
	String() string
}

// NotificationUserRecipient represents a single user Firebase message recipient.
// Once this type is used as the recipient of a notification, the user will receive a
// notification on all and only their personal devices
type NotificationUserRecipient struct {
	Address string
}

func NewNotificationUserRecipient(address string) NotificationRecipient {
	return &NotificationUserRecipient{
		Address: address,
	}
}

func (recipient *NotificationUserRecipient) GetValue() string {
	return recipient.Address
}
func (recipient *NotificationUserRecipient) String() string {
	return recipient.Address
}

// NotificationTopicRecipient represents a topic Firebase message recipient.
// Once this type if used, all the applications subscribed to the specified topic will receive the notification
type NotificationTopicRecipient struct {
	Topic string
}

func NewNotificationTopicRecipient(topic string) NotificationRecipient {
	return &NotificationTopicRecipient{
		Topic: topic,
	}
}

func (recipient *NotificationTopicRecipient) GetValue() string {
	return recipient.Topic
}
func (recipient *NotificationTopicRecipient) String() string {
	return fmt.Sprintf("topic:%s", recipient.Topic)
}

// --------------------------------------------------------------------------------------------------------------------

type NotificationData interface {
	GetType() string
	GetNotification() *messaging.Notification
	GetAdditionalData() map[string]string
}

type NotificationDataWithConfig interface {
	NotificationData
	GetAndroidConfig() *messaging.AndroidConfig
	GetAPNSConfig() *messaging.APNSConfig
	GetWebpushConfig() *messaging.WebpushConfig
}

var (
	_ NotificationDataWithConfig = StdNotificationDataWithConfig{}
)

type StdNotificationDataWithConfig struct {
	Type         string
	Data         map[string]string
	Notification *messaging.Notification
	Android      *messaging.AndroidConfig
	APNS         *messaging.APNSConfig
	Webpush      *messaging.WebpushConfig
}

func NewStdNotificationDataWithConfig(notification *messaging.Notification, data map[string]string) *StdNotificationDataWithConfig {
	if _, hasTypeField := data[NotificationTypeKey]; !hasTypeField {
		data[NotificationTypeKey] = "unknown"
	}

	if _, hasClickActionField := data[ClickActionKey]; !hasClickActionField {
		data[ClickActionKey] = ClickActionValue
	}

	return &StdNotificationDataWithConfig{
		Type:         data[NotificationTypeKey],
		Data:         data,
		Notification: notification,
	}
}

func (s StdNotificationDataWithConfig) GetType() string {
	return s.Type
}

func (s StdNotificationDataWithConfig) GetAdditionalData() map[string]string {
	return s.Data
}

func (s StdNotificationDataWithConfig) GetNotification() *messaging.Notification {
	return s.Notification
}

func (s StdNotificationDataWithConfig) GetAndroidConfig() *messaging.AndroidConfig {
	return s.Android
}

func (s StdNotificationDataWithConfig) GetAPNSConfig() *messaging.APNSConfig {
	return s.APNS
}

func (s StdNotificationDataWithConfig) GetWebpushConfig() *messaging.WebpushConfig {
	return s.Webpush
}

// --------------------------------------------------------------------------------------------------------------------

type Notification struct {
	Recipient NotificationRecipient
	Type      string
	Data      map[string]string
	Timestamp time.Time
}

func NewNotification(recipient NotificationRecipient, notificationType string, data map[string]string, timestamp time.Time) Notification {
	return Notification{
		Recipient: recipient,
		Type:      notificationType,
		Data:      data,
		Timestamp: timestamp,
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
