package handlers

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/djuno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/types"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// HandleMsgCreatePost allows to properly handle a MsgCreatePostEvent by storing inside the
// database the post that has been created with such message.
func HandleMsgCreatePost(tx types.Tx, index int, msg posts.MsgCreatePost, db db.DesmosDb) error {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Msg("Found MsgCreatePost")

	// Get the post id
	event, err := findCreationEvent(tx, index)
	if err != nil {
		return err
	}
	postID, err := findPostID(tx, event)
	if err != nil {
		return err
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

	return db.SavePost(post)
}

func findCreationEvent(tx types.Tx, index int) (sdk.StringEvent, error) {
	for _, ev := range tx.Logs[index].Events {
		if ev.Type == "post_created" {
			return ev, nil
		}
	}

	return sdk.StringEvent{}, fmt.Errorf("no post_created event found inside tx with hash %s", tx.TxHash)
}

func findPostID(tx types.Tx, event sdk.StringEvent) (posts.PostID, error) {
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

// HandleMsgEditPost allows to properly handle a MsgEditPost by updating the post inside
// the database as well.
func HandleMsgEditPost(msg posts.MsgEditPost, db db.DesmosDb) error {
	return db.EditPost(msg.PostID, msg.Message, msg.EditDate)
}

// HandleMsgAddPostReaction allows to properly handle a MsgAddPostReaction saving the
// new reaction inside the database.
func HandleMsgAddPostReaction(msg posts.MsgAddPostReaction, db db.DesmosDb) error {
	return db.SaveReaction(msg.PostID, msg.Value, msg.User)
}

// HandleMsgRemovePostReaction allows to properly handle a MsgRemovePostReaction by
// deleting the specified reaction from the database.
func HandleMsgRemovePostReaction(msg posts.MsgRemovePostReaction, db postgresql.Database) error {
	statement := `
	DELETE FROM reaction
	WHERE post_id = $1 AND owner = $2 AND value = $3;
	`

	return db.Sql.QueryRow(
		statement,
		msg.PostID, msg.User.String(), msg.Reaction,
	).Scan()
}

// HandleMsgAnswerPoll allows to properly handle a MsgAnswerPoll message by
// storing inside the database the new answer.
func HandleMsgAnswerPoll(msg posts.MsgAnswerPoll, db postgresql.Database) error {
	statement := `
	INSERT INTO user_poll_answer (poll_id, answers, user_address)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	return db.Sql.QueryRow(
		statement,
		msg.PostID, msg.UserAnswers, pq.Array(msg.Answerer),
	).Scan()
}
