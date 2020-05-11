package handlers

import (
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
	// Get the post id
	event, err := FindEventByType(tx, index, "post_created")
	if err != nil {
		return nil, err
	}
	postIDStr, err := FindAttributeByKey(tx, event, "post_id")
	if err != nil {
		return nil, err
	}
	postID, err := posts.ParsePostID(postIDStr)
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

	log.Info().Str("id", postID.String()).Str("owner", post.Creator.String()).Msg("saving post")

	// Save the post
	err = db.SavePost(post)
	if err != nil {
		return nil, err
	}

	return &post, err
}

// ____________________________________

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func HandleMsgEditPost(msg posts.MsgEditPost, db db.DesmosDb) error {
	return db.EditPost(msg.PostID, msg.Message, msg.EditDate)
}
