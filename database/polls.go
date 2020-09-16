package database

import (
	"fmt"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SavePollData allows to properly store the given poll inside the database, returning the
// id of the newly created (or updated) row inside the database itself.
// If the given poll is nil, it will not be inserted and nil will be returned as the id.
func (db DesmosDb) SavePollData(postID poststypes.PostID, poll *poststypes.PollData) error {
	// Nil data, do nothing
	if poll == nil {
		return nil
	}

	// Saving the poll data
	var pollID *uint64
	stmt := `INSERT INTO poll (post_id, question, end_date, allows_multiple_answers, allows_answer_edits)
			 VALUES ($1, $2, $3, $4, $5)
			 RETURNING id`

	err := db.Sql.QueryRow(stmt,
		postID.String(), poll.Question, poll.EndDate, poll.AllowsMultipleAnswers, poll.AllowsAnswerEdits,
	).Scan(&pollID)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO poll_answer(poll_id, answer_id, answer_text) 
			VALUES($1, $2, $3)
			ON CONFLICT ON CONSTRAINT answer_unique DO NOTHING`

	for _, answer := range poll.ProvidedAnswers {
		_, err = db.Sql.Exec(stmt, pollID, answer.ID, answer.Text)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveUserPollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SaveUserPollAnswer(postID poststypes.PostID, answer poststypes.UserAnswer) error {
	if err := db.SaveUserIfNotExisting(answer.User); err != nil {
		return err
	}

	poll, err := db.GetPollByPostID(postID)
	if err != nil {
		return err
	}
	if poll == nil {
		return fmt.Errorf("post with id %s has no poll associated to it", postID)
	}

	// Remove any existing answer to make sure that when replacing we do not get double answers
	stmt := `DELETE FROM user_poll_answer WHERE poll_id = $1 AND answerer_address = $2`
	_, err = db.Sql.Exec(stmt, poll.Id, answer.User.String())
	if err != nil {
		return err
	}

	stmt = `INSERT INTO user_poll_answer (poll_id, answer, answerer_address) 
  			VALUES ($1, $2, $3)`

	for _, answerText := range answer.Answers {
		_, err = db.Sql.Exec(stmt, poll.Id, answerText, answer.User.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPollByPostID returns the poll row associated to the post having the specified id.
// If the post with the same id has no poll associated to it, nil is returned instead.
func (db DesmosDb) GetPollByPostID(postID poststypes.PostID) (*dbtypes.PollRow, error) {
	sqlStmt := `SELECT * FROM poll WHERE post_id = $1`

	var rows []dbtypes.PollRow
	err := db.Sqlx.Select(&rows, sqlStmt, postID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}
