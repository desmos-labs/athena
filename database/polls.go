package database

import (
	poststypes "github.com/desmos-labs/desmos/x/posts/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SavePollData allows to properly store the given poll inside the database, returning the
// id of the newly created (or updated) row inside the database itself.
// If the given poll is nil, it will not be inserted and nil will be returned as the id.
func (db DesmosDb) SavePollData(postID string, poll *poststypes.PollData) error {
	// Nil data, do nothing
	if poll == nil {
		return nil
	}

	// Saving the poll data
	var pollID *uint64
	statement := `INSERT INTO poll (post_id, question, end_date, allows_multiple_answers, allows_answer_edits)
				  VALUES ($1, $2, $3, $4, $5)
				  RETURNING id`

	err := db.Sql.QueryRow(statement,
		postID, poll.Question, poll.EndDate, poll.AllowsMultipleAnswers, poll.AllowsAnswerEdits,
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

// DeletePollData allows to delete all the poll data related to the post having the given id.
func (db DesmosDb) DeletePollData(postID string) error {
	var pollID *uint64
	err := db.Sql.QueryRow(`SELECT id FROM poll WHERE post_id = $1`, postID).Scan(&pollID)
	if err != nil {
		return err
	}

	stmt := `DELETE FROM poll WHERE id = $1`
	_, err = db.Sql.Exec(stmt, pollID)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM poll_answer WHERE poll_id = $1`
	_, err = db.Sql.Exec(stmt, pollID)
	return err
}

// SaveUserPollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SaveUserPollAnswer(postID string, answer poststypes.UserAnswer) error {
	err := db.SaveUserIfNotExisting(answer.User)
	if err != nil {
		return err
	}

	statement := `
INSERT INTO user_poll_answer (poll_id, answer, answerer_address) 
VALUES ((SELECT poll_id FROM poll WHERE post_id = $1), $2, $3)`

	for _, answerText := range answer.Answers {
		_, err = db.Sql.Exec(statement, postID, answerText, answer.User)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPollByPostID returns the poll row associated to the post having the specified id.
// If the post with the same id has no poll associated to it, nil is returned instead.
func (db DesmosDb) GetPollByPostID(postID string) (*poststypes.PollData, error) {
	sqlStmt := `SELECT * FROM poll WHERE post_id = $1`

	var rows []dbtypes.PollRow
	err := db.Sqlx.Select(&rows, sqlStmt, postID)
	if err != nil {
		return nil, err
	}

	// Return nil if no poll is present
	if len(rows) == 0 {
		return nil, nil
	}

	row := rows[0]

	var answers []dbtypes.PollAnswerRow
	err = db.Sqlx.Select(&answers, `SELECT * FROM poll_answer WHERE poll_id = $1`, row.ID)
	if err != nil {
		return nil, err
	}

	return dbtypes.ConvertPollRow(row, dbtypes.ConvertPollAnswerRows(answers)), nil
}
