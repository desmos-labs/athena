package db

import (
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/lib/pq"
)

func (db DesmosDb) SavePollData(poll *posts.PollData) (pollID *uint64, err error) {
	// Nil data, do nothing
	if poll == nil {
		return nil, nil
	}

	// Saving post's poll data before post to make possible the insertion of poll_id inside it
	statement := `
		INSERT INTO poll (question, end_date, open, allows_multiple_answers, allows_answer_edits)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
		`

	err = db.Sql.QueryRow(statement,
		poll.Question, poll.EndDate, poll.Open, poll.AllowsMultipleAnswers, poll.AllowsAnswerEdits,
	).Scan(&pollID)
	if err != nil {
		return nil, err
	}

	pollQuery := `INSERT INTO poll_answer(poll_id, answer_id, answer_text) VALUES($1, $2, $3)`
	for _, answer := range poll.ProvidedAnswers {
		_, err = db.Sql.Exec(pollQuery, pollID, answer.ID, answer.Text)
		if err != nil {
			return nil, err
		}
	}

	return
}

// SavePollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SavePollAnswer(postID posts.PostID, answer posts.UserAnswer) error {
	user, err := db.SaveUserIfNotExisting(answer.User)
	if err != nil {
		return err
	}

	statement := `INSERT INTO user_poll_answer (poll_id, answers, user_id) VALUES ($1, $2, $3)`
	_, err = db.Sql.Exec(statement, postID, pq.Array(answer.Answers), user.Id)
	return err
}
