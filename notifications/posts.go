package notifications

import (
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
)

const (
	LikeReactionValue = ":heart:"

	NotificationTypeKey = "type"
	TypeComment         = "comment"
	TypeReaction        = "reaction"
	TypeLike            = "like"

	NotificationActionKey = "action"
	ActionOpenPost        = "open_post"

	PostIDKey      = "post_id"
	PostMessageKey = "post_message"
	PostCreatorKey = "post_creator"

	PostReactionValueKey = "post_reaction_value"
	PostReactionOwnerKey = "post_reaction_owner"

	PostLikeUserKey = "like_user"
)

// SendPostNotifications takes the given post and, upon having performed the necessary checks, sends
// a push notification to the people that might somehow be interested into the creation of the post.
// For example, if the post is a comment to another post, the creator of the latter will be notified
// that a new comment has been added.
func SendPostNotifications(post posts.Post, db db.DesmosDb) error {
	// Get the post parent
	parent, err := db.GetPostByID(post.ParentID)
	if err != nil {
		return err
	}

	// Send the notification as it's a comment
	if err := sendCommentNotification(post, parent); err != nil {
		return err
	}

	// TODO: Send mention notification
	// TODO: Send tag notification

	return nil
}

func sendCommentNotification(post posts.Post, parent *posts.Post) error {
	// Not a comment, skip
	if parent == nil {
		return nil
	}

	// Post and comment creators are the same, just return
	if post.Creator.Equals(parent.Creator) {
		return nil
	}

	// Build the notification
	notification := messaging.Notification{
		Title: "Someone commented one of your posts! 💬",
		Body:  fmt.Sprintf("%s commented on your post: %s", post.Creator.String(), post.Message),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeComment,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:      post.PostID.String(),
		PostMessageKey: post.Message,
		PostCreatorKey: post.Creator.String(),
	}

	// Send a notification to the original post owner
	return SendNotification(parent.Creator.String(), &notification, data)
}

// SendReactionNotifications takes the given reaction (which has been added to the post having the given id)
// and sends out push notifications to all the users that might be interested in the reaction creation event.
// For example, a push notification is send to the user that has created the post.
func SendReactionNotifications(postID posts.PostID, reaction posts.Reaction, db db.DesmosDb) error {
	post, err := db.GetPostByID(postID)
	if err != nil {
		return err
	}

	// The post creator and the reaction owner are the same person, just return
	if post.Creator.Equals(reaction.Owner) {
		return nil
	}

	if reaction.Value == LikeReactionValue {
		return sendLikeNotification(post, reaction)
	}
	return sendGenericReactionNotification(post, reaction)
}

// sendGenericReactionNotification allows to send a notification for a generic given reaction
// that has been added to the specified post
func sendGenericReactionNotification(post *posts.Post, reaction posts.Reaction) error {
	// Build the notification
	notification := messaging.Notification{
		Title: "Someone added a new reaction! 🎉",
		Body:  fmt.Sprintf("%s added a new reaction to your post: %s", reaction.Owner.String(), reaction.Value),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeReaction,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:            post.PostID.String(),
		PostReactionValueKey: reaction.Value,
		PostReactionOwnerKey: reaction.Owner.String(),
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator.String(), &notification, data)
}

// sendLikeNotification sends a push notification telling that a like has been added to the given post
func sendLikeNotification(post *posts.Post, reaction posts.Reaction) error {
	// Build the notification
	notification := messaging.Notification{
		Title: "Someone like one of your posts ❤️",
		Body:  fmt.Sprintf("%s like your post!", reaction.Owner.String()),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeLike,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:       post.PostID.String(),
		PostLikeUserKey: reaction.Owner.String(),
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator.String(), &notification, data)
}
