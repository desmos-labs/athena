package builder

import (
	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"

	"github.com/desmos-labs/djuno/v2/types"
	notificationscontext "github.com/desmos-labs/djuno/v2/x/notifications/context"
)

// -------------------------------------------------------------------------------------------------------------------

type NotificationsBuilderCreator func(context notificationscontext.Context) NotificationsBuilder

// NotificationsBuilder contains all the notifications builders
type NotificationsBuilder interface {
	Posts() PostsNotificationsBuilder
	Reactions() ReactionsNotificationsBuilder
	Relationships() RelationshipsNotificationsBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type PostNotificationBuilder = func(originalPost types.Post, post types.Post) types.NotificationData

type MentionNotificationBuilder = func(post types.Post, mention poststypes.TextTag) types.NotificationData

// PostsNotificationsBuilder contains all the notifications builders for the posts module
type PostsNotificationsBuilder interface {
	Comment() PostNotificationBuilder
	Reply() PostNotificationBuilder
	Repost() PostNotificationBuilder
	Quote() PostNotificationBuilder
	Mention() MentionNotificationBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type ReactionNotificationBuilder = func(post types.Post, reaction types.Reaction) types.NotificationData

// ReactionsNotificationsBuilder contains all the notifications builders for the reactions module
type ReactionsNotificationsBuilder interface {
	Reaction() ReactionNotificationBuilder
}

// -------------------------------------------------------------------------------------------------------------------

type RelationshipNotificationBuilder = func(relationship types.Relationship) types.NotificationData

// RelationshipsNotificationsBuilder contains all the notifications builders for the relationships module
type RelationshipsNotificationsBuilder interface {
	Relationship() RelationshipNotificationBuilder
}
