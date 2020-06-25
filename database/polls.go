package database

import (
	"fmt"

	"github.com/desmos-labs/desmos/x/posts"
	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SavePollData allows to properly store the given poll inside the database, returning the
// id of the newly created (or updated) row inside the database itself.
// If the given poll is nil, it will not be inserted and nil will be returned as the id.
func (db DesmosDb) SavePollData(postID posts.PostID, poll *posts.PollData) error {
	// Nil data, do nothing
	if poll == nil {
		return nil
	}

	// Saving the poll data
	var pollID *uint64
	statement := `INSERT INTO poll (post_id, question, end_date, open, allows_multiple_answers, allows_answer_edits)
				  VALUES ($1, $2, $3, $4, $5, $6)
				  RETURNING id`

	err := db.Sql.QueryRow(statement,
		postID.String(), poll.Question, poll.EndDate, poll.Open, poll.AllowsMultipleAnswers, poll.AllowsAnswerEdits,
	).Scan(&pollID)
	if err != nil {
		return err
	}

	pollQuery := `INSERT INTO poll_answer(poll_id, answer_id, answer_text) VALUES($1, $2, $3)`
	for _, answer := range poll.ProvidedAnswers {
		_, err = db.Sql.Exec(pollQuery, pollID, answer.ID, answer.Text)
		if err != nil {
			return err
		}
	}

	return nil
}

// SavePollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SavePollAnswer(postID posts.PostID, answer posts.UserAnswer) error {
	_, err := db.SaveUserIfNotExisting(answer.User)
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

	statement := `INSERT INTO user_poll_answer (poll_id, answer, answerer_address) VALUES ($1, $2, $3)`
	for _, answerText := range answer.Answers {
		_, err = db.Sql.Exec(statement, poll.Id, answerText, answer.User.String())
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPollByPostID returns the poll row associated to the post having the specified id.
// If the post with the same id has no poll associated to it, nil is returned instead.
func (db DesmosDb) GetPollByPostID(postID posts.PostID) (*dbtypes.PollRow, error) {
	sqlStmt := `SELECT * FROM poll WHERE post_id = ?`

	var rows []dbtypes.PollRow
	err := db.sqlx.Select(&rows, sqlStmt, postID)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}
