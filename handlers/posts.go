package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/djuno/notifications"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
)

// ____________________________________

// HandleMsgCreatePost allows to properly handle the given msg present inside the specified tx at the specific
// index. It creates a new Post object from it, stores it inside the database and later sends out any
// push notification using Firebase Cloud Messaging.
func HandleMsgCreatePost(tx types.Tx, index int, msg posts.MsgCreatePost, db db.DesmosDb) error {
	post, err := CreateAndStorePostFromMsgCreatePost(tx, index, msg, db)
	if err != nil {
		return err
	}

	return notifications.SendPostNotifications(*post, db)
}

// CreateAndStorePostFromMsgCreatePost allows to properly handle a MsgCreatePostEvent by storing inside the
// database the post that has been created with such message.
// After the post has been saved, it is returned for other uses.
func CreateAndStorePostFromMsgCreatePost(tx types.Tx, index int, msg posts.MsgCreatePost, db db.DesmosDb) (*posts.Post, error) {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Msg("Found MsgCreatePost")

	// Get the post id
	event, err := FindCreationEvent(tx, index)
	if err != nil {
		return nil, err
	}
	postID, err := FindPostID(tx, event)
	if err != nil {
		return nil, err
	}

	// Create the post
	post := posts.NewPost(postID, msg.ParentID, msg.Message, msg.AllowsComments,
		msg.Subspace, msg.OptionalData, msg.CreationDate, msg.Creator)

	if msg.Medias != nil {
		post = post.WithMedias(msg.Medias)
	}

	if msg.PollData != nil {
		post = post.WithPollData(*msg.PollData)
	}

	// Save the post
	err = db.SavePost(post)
	if err != nil {
		return nil, err
	}

	return &post, err
}

// FindCreationEvent searches inside the given tx events for the message having the specified index, in order
// to find the event related to a post creation, and returns it.
// If no such event is found, returns an error instead.
func FindCreationEvent(tx types.Tx, index int) (sdk.StringEvent, error) {
	for _, ev := range tx.Logs[index].Events {
		if ev.Type == "post_created" {
			return ev, nil
		}
	}

	return sdk.StringEvent{}, fmt.Errorf("no post_created event found inside tx with hash %s", tx.TxHash)
}

// FindPostID searches inside the specified event of the given tx to find the newly generated id of a post.
// If the specified event does not contain a new post PostID, returns an error instead.
func FindPostID(tx types.Tx, event sdk.StringEvent) (posts.PostID, error) {
	for _, attr := range event.Attributes {
		if attr.Key == "post_id" {
			postID, err := posts.ParsePostID(attr.Value)
			if err != nil {
				return posts.PostID(0), err
			}

			return postID, nil
		}
	}

	return posts.PostID(0), fmt.Errorf("no event with attribute post_id found inside tx with hash %s", tx.TxHash)
}

// ____________________________________

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func HandleMsgEditPost(msg posts.MsgEditPost, db db.DesmosDb) error {
	return db.EditPost(msg.PostID, msg.Message, msg.EditDate)
}

// ____________________________________

// HandleMsgAddPostReaction allows to properly handle a MsgAddPostReaction by storing the newly created
// reaction inside the database and sending out push notifications to whoever might be interested in this event.
func HandleMsgAddPostReaction(msg posts.MsgAddPostReaction, db db.DesmosDb) error {
	reaction, err := db.SaveReaction(msg.PostID, posts.NewReaction(msg.Value, msg.User))
	if err != nil {
		return err
	}

	return notifications.SendReactionNotifications(msg.PostID, *reaction, db)
}

// ____________________________________

// HandleMsgRemovePostReaction allows to properly handle a MsgRemovePostReaction by
// deleting the specified reaction from the database.
func HandleMsgRemovePostReaction(msg posts.MsgRemovePostReaction, db db.DesmosDb) error {
	return db.RemoveReaction(msg.PostID, msg.Reaction, msg.User)
}

// ____________________________________

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func HandleMsgAnswerPoll(msg posts.MsgAnswerPoll, db db.DesmosDb) error {
	return db.SavePollAnswer(msg.PostID, posts.NewUserAnswer(msg.UserAnswers, msg.Answerer))
}
