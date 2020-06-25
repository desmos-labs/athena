package notifications

import (
	"fmt"
	"regexp"
	"strings"

	"firebase.google.com/go/messaging"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
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
	mentionRegEx = regexp.MustCompile(`\s([@][a-zA-Z_1-9]+)`)
)

// SendPostNotifications takes the given post and, upon having performed the necessary checks, sends
// a push notification to the people that might somehow be interested into the creation of the post.
// For example, if the post is a comment to another post, the creator of the latter will be notified
// that a new comment has been added.
func SendPostNotifications(post posts.Post, db database.DesmosDb) error {
	// Get the post parent
	parent, err := db.GetPostByID(post.ParentID)
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

// sendMentionNotifications sends everyone who is tagged inside the given post message a notification.
// If the given post is a comment to another post, the notification will not be sent to the user that has
// created the post to which this post is a comment. He will already receive the comment notification,
// so we need to avoid double notifications
func sendMentionNotifications(parent *posts.Post, post posts.Post) error {

	var originalPoster sdk.AccAddress
	if parent != nil {
		originalPoster = parent.Creator
	}

	mentions, err := GetPostMentions(post)
	if err != nil {
		return err
	}

	for _, address := range mentions {
		// No notification to the original poster
		if len(originalPoster) != 0 && address.Equals(originalPoster) {
			continue
		}

		// No notification to the post creator if he has mentioned himself
		if address.Equals(post.Creator) {
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
func GetPostMentions(post posts.Post) ([]sdk.AccAddress, error) {
	mentions := mentionRegEx.FindAllString(post.Message, -1)

	addresses := make([]sdk.AccAddress, len(mentions))
	for index, mention := range mentions {
		addressString := strings.Trim("@", strings.TrimSpace(mention))
		address, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return nil, err
		}
		addresses[index] = address
	}

	return addresses, nil
}

// sendMentionNotification sends a single notification to the given telling him that he's been mentioned
// inside the given post.
func sendMentionNotification(post posts.Post, user sdk.AccAddress) error {
	// Get the mentions
	notification := messaging.Notification{
		Title: "You've been mentioned inside a post",
		Body:  fmt.Sprintf("%s has mentioned you inside a post: %s", post.Creator, post.Message),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeMention,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:          post.PostID.String(),
		PostMentionUserKey: post.Creator.String(),
		PostMentionTextKey: post.Message,
	}

	return SendNotification(user.String(), &notification, data)
}

// SendReactionNotifications takes the given reaction (which has been added to the post having the given id)
// and sends out push notifications to all the users that might be interested in the reaction creation event.
// For example, a push notification is send to the user that has created the post.
func SendReactionNotifications(postID posts.PostID, reaction *posts.PostReaction, db database.DesmosDb) error {
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
func sendGenericReactionNotification(post *posts.Post, reaction *posts.PostReaction) error {
	// Build the notification
	notification := messaging.Notification{
		Title: "Someone added a new reaction to one of your posts 🎉",
		Body:  fmt.Sprintf("%s added a new reaction to your post: %s", reaction.Owner.String(), reaction.Value),
	}
	data := map[string]string{
		NotificationTypeKey:   TypeReaction,
		NotificationActionKey: ActionOpenPost,

		PostIDKey:                post.PostID.String(),
		PostReactionValueKey:     reaction.Value,
		PostReactionShortCodeKey: reaction.Shortcode,
		PostReactionOwnerKey:     reaction.Owner.String(),
	}

	// Send a notification to the post creator
	return SendNotification(post.Creator.String(), &notification, data)
}

// sendLikeNotification sends a push notification telling that a like has been added to the given post
func sendLikeNotification(post *posts.Post, reaction *posts.PostReaction) error {
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
