package utils

import (
	"regexp"
	"strings"
	"time"

	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/desmos-labs/djuno/types"
)

var (
	mentionRegEx = regexp.MustCompile(`\s([@][a-zA-Z0-9]+)`)
)

// GetPostMentions returns the list of all the addresses that have been mentioned inside a post.
// If no mentions are present, returns nil instead.
func GetPostMentions(post *types.Post) ([]string, error) {
	mentions := mentionRegEx.FindAllString(post.Message, -1)

	addresses := make([]string, len(mentions))
	for index, mention := range mentions {
		addresses[index] = strings.Trim(strings.TrimSpace(mention), "@")
	}

	return addresses, nil
}

// GetPostFromMsgCreatePost creates
func GetPostFromMsgCreatePost(tx *juno.Tx, index int, msg *poststypes.MsgCreatePost) (*types.Post, error) {
	event, err := tx.FindEventByType(index, poststypes.EventTypePostCreated)
	if err != nil {
		return nil, err
	}

	// Get the post id
	postID, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return nil, err
	}

	// Get the creation time
	creationTimeStr, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostCreationTime)
	if err != nil {
		return nil, err
	}
	creationTime, err := time.Parse(time.RFC3339, creationTimeStr)
	if err != nil {
		return nil, err
	}

	// Create the post
	return types.NewPost(
		poststypes.NewPost(
			postID,
			msg.ParentID,
			msg.Message,
			msg.CommentsState,
			msg.Subspace,
			msg.AdditionalAttributes,
			msg.Attachments,
			msg.PollData,
			creationTime,
			time.Time{},
			msg.Creator,
		),
		tx.Height,
	), nil
}

// GetReactionFromTxEvent creates a new PostReaction object from the event having the given type and associated
// to the message having the given inside the inside the given tx.
func GetReactionFromTxEvent(tx *juno.Tx, index int, eventType string) (string, poststypes.PostReaction, error) {
	event, err := tx.FindEventByType(index, eventType)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	postID, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostID)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	user, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionOwner)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	value, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyPostReactionValue)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	shortCode, err := tx.FindAttributeByKey(event, poststypes.AttributeKeyReactionShortCode)
	if err != nil {
		return "", poststypes.PostReaction{}, err
	}

	return postID, poststypes.NewPostReaction(postID, shortCode, value, user), nil
}
