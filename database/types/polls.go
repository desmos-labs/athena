package types

import (
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
)

// PollRow represents a single PostgreSQL row containing the details of a poll
type PollRow struct {
	PostID                string    `db:"post_id"`
	ID                    uint64    `db:"id"`
	Question              string    `db:"question"`
	EndDate               time.Time `db:"end_date"`
	AllowsMultipleAnswers bool      `db:"allows_multiple_answers"`
	AllowsAnswerEdits     bool      `db:"allows_answer_edits"`
}

// Equal tells whether r and s contain the same data
func (r PollRow) Equal(s PollRow) bool {
	return r.PostID == s.PostID &&
		r.ID == s.ID &&
		r.Question == s.Question &&
		r.EndDate.Equal(s.EndDate) &&
		r.AllowsMultipleAnswers == s.AllowsMultipleAnswers &&
		r.AllowsAnswerEdits == s.AllowsAnswerEdits
}

// ConvertPollRow converts the given row and answers into a proper PollData object
func ConvertPollRow(row PollRow, answers []poststypes.PollAnswer) *poststypes.PollData {
	return poststypes.NewPollData(
		row.Question,
		row.EndDate,
		answers,
		row.AllowsMultipleAnswers,
		row.AllowsAnswerEdits,
	)
}

// ---------------------------------------------------------------------------------------------------

// PollAnswerRow represents a single row of the poll_answer table
type PollAnswerRow struct {
	PollID   int    `db:"poll_id"`
	AnswerID string `db:"answer_id"`
	Text     string `db:"answer_text"`
}

// ConvertPollAnswerRows converts the given rows into PollAnswer objects
func ConvertPollAnswerRows(rows []PollAnswerRow) []poststypes.PollAnswer {
	answers := make([]poststypes.PollAnswer, len(rows))
	for index, rows := range rows {
		answers[index] = poststypes.NewPollAnswer(rows.AnswerID, rows.Text)
	}
	return answers
}

// ---------------------------------------------------------------------------------------------------

// UserPollAnswerRow represents a single row of the user_poll_answer table
type UserPollAnswerRow struct {
	PollID          int    `db:"poll_id"`
	Answer          string `db:"answer"`
	AnswererAddress string `db:"answerer_address"`
}

// Equal tells whether r and s contain the same data
func (r UserPollAnswerRow) Equal(s UserPollAnswerRow) bool {
	return r.PollID == s.PollID &&
		r.Answer == s.Answer &&
		r.AnswererAddress == s.AnswererAddress
}
