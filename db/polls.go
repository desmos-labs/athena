package db

import (
	"fmt"
	"time"

	"github.com/desmos-labs/desmos/x/posts"
)

type PollRow struct {
	Id                    uint64    `db:"id"`
	Question              string    `db:"question"`
	EndDate               time.Time `db:"end_date"`
	Open                  bool      `db:"open"`
	AllowsMultipleAnswers bool      `db:"allows_multiple_answers"`
	AllowsAnswerEdits     bool      `db:"allows_answer_edits"`
}

func (db DesmosDb) GetPollByPostID(postID posts.PostID) (*PollRow, error) {
	sqlStmt := `SELECT * from poll WHERE id = (SELECT post.poll_id FROM post WHERE post.id = $1)`

	var rows []PollRow
	err := db.sqlx.Select(&rows, sqlStmt, postID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}

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

	err = db.Sql.QueryRow(
		statement,
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

	poll, err := db.GetPollByPostID(postID)
	if err != nil {
		return err
	}
	if poll == nil {
		return fmt.Errorf("post with id %s has no poll associated to it", postID)
	}

	statement := `INSERT INTO user_poll_answer (poll_id, answer, user_id) VALUES ($1, $2, $3)`
	for _, answer := range answer.Answers {
		_, err = db.Sql.Exec(statement, poll.Id, answer, user.Id)
		if err != nil {
			return err
		}
	}

	return nil
}
