package handlers

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/lib/pq"
)

func handleMsgAnswerPoll(msg posts.MsgAnswerPoll, db postgresql.Database) error {
	var id uint64

	addPollAnswersSqlStatement := `
	INSERT INTO user_poll_answer (poll_id, answers, user_address)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	return db.Sql.QueryRow(
		addPollAnswersSqlStatement,
		msg.PostID,
		pq.Array(msg.UserAnswers),
		msg.Answerer.String(),
	).Scan(&id)
}
