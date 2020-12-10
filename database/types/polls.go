package types

import (
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
)

// PollRow represents a single PostgreSQL row containing the details of a poll
type PollRow struct {
	PostID                string    `db:"post_id"`
	Id                    uint64    `db:"id"`
	Question              string    `db:"question"`
	EndDate               time.Time `db:"end_date"`
	AllowsMultipleAnswers bool      `db:"allows_multiple_answers"`
	AllowsAnswerEdits     bool      `db:"allows_answer_edits"`
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
