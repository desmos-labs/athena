package notifications

import (
	"fmt"
	"regexp"
	"strings"

	"firebase.google.com/go/messaging"
	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"

	"github.com/desmos-labs/djuno/database"
)

const (
	LikeReactionValue = ":heart:"

	NotificationTypeKey = "type"
	TypeComment         = "comment"
	TypeReaction        = "reaction"
	TypeLike            = "like"
	TypeMention         = "mention"

	NotificationActionKey = "action"
	ActionOpenPost        = "open_post"

	PostIDKey      = "post_id"
	PostMessageKey = "post_message"
	PostCreatorKey = "post_creator"

	PostReactionValueKey     = "post_reaction_value"
	PostReactionShortCodeKey = "post_reaction_shortcode"
	PostReactionOwnerKey     = "post_reaction_owner"

	PostLikeUserKey = "like_user"

	PostMentionUserKey = "mention_user"
	PostMentionTextKey = "mention_text"
)

var (
	mentionRegEx = regexp.MustCompile(`\s([@][a-zA-Z0-9]+)`)
)

// SendPostNotifications takes the given post and, upon having performed the necessary checks, sends
// a push notification to the people that might somehow be interested into the creation of the post.
// For example, if the post is a comment to another post, the creator of the latter will be notified
// that a new comment has been added.
func SendPostNotifications(post poststypes.Post, db *database.DesmosDb) error {
	// Get the post parent
	parent, err := db.GetPostByID(post.ParentId)
	if err != nil {
		return err
	}

	// Send the notification as it's a comment
	if err := sendCommentNotification(post, parent); err != nil {
		return err
	}

	// Send the mentions notifications
	if err := sendMentionNotifications(parent, post); err != nil {
		return err
	}

	// TODO: Send tag notification

	return nil
}

// sendCommentNotification sends the creator of the parent post a notification telling him
// that the given post has been added as a comment to it's original post.
func sendCommentNotification(post poststypes.Post, parent *poststypes.Post) error {
	// Not a comment, skip
	if parent == nil {
		return nil
	}

	// Post and comment creators are the same, just return
	if post.Creator == parent.Creator {
		return nil
	}

	// Build the notification
	notification := messaging.Notification{
		Title: "Someone commented one of your posts! 💬",
		Body:  fmt.Sprintf("%s commented on your post: %s", post.Creator, post.Message),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeComment,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:      post.PostId,
		PostMessageKey: post.Message,
		PostCreatorKey: post.Creator,
	}

	// Send a notification to the original post owner
	return SendNotification(parent.Creator, &notification, data)
}

// sendMentionNotifications sends everyone who is tagged inside the given post message a notification.
// If the given post is a comment to another post, the notification will not be sent to the user that has
// created the post to which this post is a comment. He will already receive the comment notification,
// so we need to avoid double notifications
func sendMentionNotifications(parent *poststypes.Post, post poststypes.Post) error {

	var originalPoster string
	if parent != nil {
		originalPoster = parent.Creator
	}

	mentions, err := GetPostMentions(post)
	if err != nil {
		return err
	}

	for _, address := range mentions {
		// No notification to the original poster
		if len(originalPoster) != 0 && address == originalPoster {
			continue
		}

		// No notification to the post creator if he has mentioned himself
		if address == post.Creator {
			continue
		}

		err := sendMentionNotification(post, address)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetPostMentions returns the list of all the addresses that have been mentioned inside a post.
// If no mentions are present, returns nil instead.
func GetPostMentions(post poststypes.Post) ([]string, error) {
	mentions := mentionRegEx.FindAllString(post.Message, -1)

	addresses := make([]string, len(mentions))
	for index, mention := range mentions {
		addresses[index] = strings.Trim(strings.TrimSpace(mention), "@")
	}

	return addresses, nil
}

// sendMentionNotification sends a single notification to the given telling him that he's been mentioned
// inside the given post.
func sendMentionNotification(post poststypes.Post, user string) error {
	// Get the mentions
	notification := messaging.Notification{
		Title: "You've been mentioned inside a post",
		Body:  fmt.Sprintf("%s has mentioned you inside a post: %s", post.Creator, post.Message),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeMention,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:          post.PostId,
		PostMentionUserKey: post.Creator,
		PostMentionTextKey: post.Message,
	}

	return SendNotification(user, &notification, data)
}

// SendReactionNotifications takes the given reaction (which has been added to the post having the given id)
// and sends out push notifications to all the users that might be interested in the reaction creation event.
// For example, a push notification is send to the user that has created the post.
func SendReactionNotifications(postID string, reaction poststypes.PostReaction, db *database.DesmosDb) error {
	post, err := db.GetPostByID(postID)
	if err != nil {
		return err
	}

	// The post creator and the reaction owner are the same person, just return
	if post.Creator == reaction.Owner {
		return nil
	}

	if reaction.Value == LikeReactionValue {
		return sendLikeNotification(post, reaction)
	}
	return sendGenericReactionNotification(post, reaction)
}

// sendGenericReactionNotification allows to send a notification for a generic given reaction
// that has been added to the specified post
func sendGenericReactionNotification(post *poststypes.Post, reaction poststypes.PostReaction) error {
	// Build the notification
	notification := messaging.Notification{
		Title: "Someone added a new reaction to one of your posts 🎉",
		Body:  fmt.Sprintf("%s added a new reaction to your post: %s", reaction.Owner, reaction.Value),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeReaction,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:                post.PostId,
		PostReactionValueKey:     reaction.Value,
		PostReactionShortCodeKey: reaction.ShortCode,
		PostReactionOwnerKey:     reaction.Owner,
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator, &notification, data)
}

// sendLikeNotification sends a push notification telling that a like has been added to the given post
func sendLikeNotification(post *poststypes.Post, reaction poststypes.PostReaction) error {
	// Build the notification
	notification := messaging.Notification{
		Title: "Someone like one of your posts ❤️",
		Body:  fmt.Sprintf("%s like your post!", reaction.Owner),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeLike,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:       post.PostId,
		PostLikeUserKey: reaction.Owner,
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator, &notification, data)
}
