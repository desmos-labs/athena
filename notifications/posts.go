package notifications

import (
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
)

const (
	PostIDKey = "post_id"

	PostMessageKey = "post_message"
	PostCreatorKey = "post_creator"

	PostReactionValueKey = "post_reaction_value"
	PostReactionOwnerKey = "post_reaction_owner"
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

	// No parent, simply return
	if parent == nil {
		return nil
	}

	// Post and comment creators are the same, just return
	if post.Creator.Equals(parent.Creator) {
		return nil
	}

	// Build the notification
	// TODO: Improve the messages here
	notification := messaging.Notification{
		Title: "New post!",
		Body:  "A new post has been created!",
	}
	data := map[string]string{
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

	// Build the notification
	notification := messaging.Notification{
		Title: "Someone added a new reaction! 🎉",
		Body:  fmt.Sprintf("%s added a new reaction to your post: %s", reaction.Owner.String(), reaction.Value),
	}
	data := map[string]string{
		PostIDKey:            postID.String(),
		PostReactionValueKey: reaction.Value,
		PostReactionOwnerKey: reaction.Owner.String(),
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator.String(), &notification, data)
}
