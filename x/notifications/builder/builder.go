package builder

import (
	"firebase.google.com/go/v4/messaging"
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// NotificationData contains the notification data returned by a generic builder
type NotificationData struct {
	Notification *messaging.Notification
	Data         map[string]string
}

// UtilityModule represents a module that contains utility method within it
type UtilityModule interface {
	GetDisplayName(userAddress string) string
}

// -------------------------------------------------------------------------------------------------------------------

type NotificationsBuilderCreator = func(module UtilityModule) NotificationsBuilder

// NotificationsBuilder contains all the notifications builders
type NotificationsBuilder interface {
	Posts() PostsNotificationsBuilder
	Reactions() ReactionsNotificationsBuilder
	Relationships() RelationshipsNotificationsBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type PostNotificationBuilder = func(originalPost types.Post, post types.Post) *NotificationData

type MentionNotificationBuilder = func(post types.Post, mention poststypes.TextTag) *NotificationData

// PostsNotificationsBuilder contains all the notifications builders for the posts module
type PostsNotificationsBuilder interface {
	ConversationReply() PostNotificationBuilder
	Comment() PostNotificationBuilder
	Repost() PostNotificationBuilder
	Quote() PostNotificationBuilder
	Mention() MentionNotificationBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type ReactionNotificationBuilder = func(post types.Post, reaction types.Reaction) *NotificationData

// ReactionsNotificationsBuilder contains all the notifications builders for the reactions module
type ReactionsNotificationsBuilder interface {
	Reaction() ReactionNotificationBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type RelationshipNotificationBuilder = func(relationship types.Relationship) *NotificationData

// RelationshipsNotificationsBuilder contains all the notifications builders for the relationships module
type RelationshipsNotificationsBuilder interface {
	Relationship() RelationshipNotificationBuilder
}
