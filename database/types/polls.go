package types

import (
	"time"
)

// PollRow represents a single PostgreSQL row containing the details of a poll
type PollRow struct {
	Id                    uint64    `db:"id"`
	PostID                string    `db:"post_id"`
	Question              string    `db:"question"`
	EndDate               time.Time `db:"end_date"`
	Open                  bool      `db:"open"`
	AllowsMultipleAnswers bool      `db:"allows_multiple_answers"`
	AllowsAnswerEdits     bool      `db:"allows_answer_edits"`
}

func (row PollRow) Equal(other PollRow) bool {
	return row.Id == other.Id &&
		row.PostID == other.PostID &&
		row.Question == other.Question &&
		row.EndDate.Equal(other.EndDate) &&
		row.Open == other.Open &&
		row.AllowsMultipleAnswers == other.AllowsMultipleAnswers &&
		row.AllowsAnswerEdits == other.AllowsAnswerEdits
}

// ________________________________________________

type UserPollAnswerRow struct {
	PollId          int64  `db:"poll_id"`
	Answer          int64  `db:"answer"`
	AnswererAddress string `db:"answerer_address"`
}

func (row UserPollAnswerRow) Equal(other UserPollAnswerRow) bool {
	return row.PollId == other.PollId &&
		row.Answer == other.Answer &&
		row.AnswererAddress == other.AnswererAddress
}
