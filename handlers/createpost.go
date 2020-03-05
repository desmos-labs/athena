package handlers

import (
	"encoding/json"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
	"strconv"
)

func handleMsgCreatePost(tx types.Tx, index int, msg posts.MsgCreatePost, db postgresql.Database) error {
	log.Info().Str("tx_hash", tx.TxHash).Int("msg_index", index).Msg("Found MsgCreatePost")

	var postID uint64

	// Get the post id
	// TODO: test with multiple MsgCreatePost
	for _, ev := range tx.Logs[index].Events {
		for _, attr := range ev.Attributes {
			if attr.Key == "post_id" {
				postID, _ = strconv.ParseUint(attr.Value, 10, 64)
				break
			}
		}
	}

	if err := savePost(postID, msg, db); err != nil {
		return err
	}

	return nil
}

// handleMsgCreatePost handles a MsgCreatePost and saves the post inside the database
func savePost(postID uint64, msg posts.MsgCreatePost, db postgresql.Database) error {
	var id *uint64

	// Saving post's poll data before post to make possible the insertion of poll_id inside it

	pollDataSqlStatement := `
	INSERT INTO poll (question, end_date, open, allows_multiple_answers, allows_answer_edits)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id;
	`

	if msg.PollData != nil {
		err := db.Sql.QueryRow(
			pollDataSqlStatement,
			msg.PollData.Question,
			msg.PollData.EndDate,
			msg.PollData.Open,
			msg.PollData.AllowsMultipleAnswers,
			msg.PollData.AllowsAnswerEdits,
		).Scan(&id)

		if err != nil {
			return err
		}

		addPollAnswersSqlStatement := `
		INSERT INTO poll_answer(poll_id, answer_id, answer_text)
		VALUES($1, $2, $3)
		RETURNING id;
		`

		for _, answer := range msg.PollData.ProvidedAnswers {
			_, err := db.Sql.Exec(
				addPollAnswersSqlStatement,
				id,
				answer.ID,
				answer.Text,
			)

			if err != nil {
				return err
			}
		}

	}

	// Saving Post

	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator, poll_id, optional_data)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;
    `

	// todo look how this is inserted in DB
	jsonB, _ := json.Marshal(msg.OptionalData)

	err := db.Sql.QueryRow(
		postSqlStatement,
		postID,
		msg.ParentID,
		msg.Message,
		msg.CreationDate,
		msg.CreationDate,
		msg.AllowsComments,
		msg.Subspace,
		msg.Creator.String(),
		id,
		string(jsonB),
	).Scan(&id)

	if err != nil {
		return err
	}
	// Saving post's medias
	mediasSqlStatement := `
	INSERT INTO media (post_id, uri, mime_type)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	for _, media := range msg.Medias {
		err = db.Sql.QueryRow(
			mediasSqlStatement,
			postID,
			media.URI,
			media.MimeType,
		).Scan(&id)

		if err != nil {
			return err
		}

	}

	return nil
}
