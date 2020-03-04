package handlers

import (
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
	var id uint64

	// Saving Post

	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id;
    `

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
	).Scan(&id)

	if err != nil {
		return err
	}

	// Saving post's optional data
	optionalDataSqlStatement := `
	INSERT INTO optional_data (id, key, value)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	for key, value := range msg.OptionalData {
		err = db.Sql.QueryRow(
			optionalDataSqlStatement,
			postID,
			key,
			value,
		).Scan(&id)

		if err != nil {
			return err
		}
	}

	// Saving post's medias
	mediasSqlStatement := `
	INSERT INTO media (id, uri, mime_type)
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

	// Saving post's poll data

	pollDataSqlStatement := `
	INSERT INTO poll_data (id, question, end_date, open, allows_multiple_answers, allows_answer_edits)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`

	if msg.PollData != nil {
		err = db.Sql.QueryRow(
			pollDataSqlStatement,
			postID,
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
		INSERT INTO poll_data_answers(id, answer_id, answer_text)
		VALUES($1, $2, $3)
		RETURNING id;
		`
		for _, answer := range msg.PollData.ProvidedAnswers {
			err := db.Sql.QueryRow(
				addPollAnswersSqlStatement,
				postID,
				answer.ID,
				answer.Text,
			).Scan(&id)

			if err != nil {
				return err
			}
		}

	}

	return nil
}
