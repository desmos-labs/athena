package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

// handleMsgCreatePost handles a MsgCreatePost and saves the post inside the database
func handleMsgCreatePost(postID uint64, msg posts.MsgCreatePost, db postgresql.Database) error {
	var id uint64

	// Saving Post

	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, allows_comments, subspace, creator)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
    `

	err := db.Sql.QueryRow(
		postSqlStatement,
		postID,
		msg.ParentID,
		msg.Message,
		msg.CreationDate,
		msg.AllowsComments,
		msg.Subspace,
		msg.Creator,
	).Scan(&id)

	if err != nil {
		return err
	}

	// Saving post's optional data
	optionalDataSqlStatement := `
	INSERT INTO optional_data (id, key, value)
	VALUES ($1, $2, $3)
	RETURNING id
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
	INSERT INTO media (id, post_medias)
	VALUES ($1, $2)
	`

	err = db.Sql.QueryRow(
		mediasSqlStatement,
		postID,
		msg.Medias,
	).Scan(&id)

	if err != nil {
		return err
	}

	// Saving post's poll data

	pollDataSqlStatement := `
	INSERT INTO poll_data (id, question, provided_answers, end_date, open, allows_multiple_answers, allows_answer_edits)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	err = db.Sql.QueryRow(
		pollDataSqlStatement,
		postID,
		msg.PollData.Question,
		msg.PollData.ProvidedAnswers,
		msg.PollData.EndDate,
		msg.PollData.Open,
		msg.PollData.AllowsMultipleAnswers,
		msg.PollData.AllowsAnswerEdits,
	).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}
