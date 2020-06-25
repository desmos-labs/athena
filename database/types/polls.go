package types

import (
	"time"
)

// PollRow represents a single PostgreSQL row containing the details of a poll
type PollRow struct {
	Id                    uint64    `db:"id"`
	Question              string    `db:"question"`
	EndDate               time.Time `db:"end_date"`
	Open                  bool      `db:"open"`
	AllowsMultipleAnswers bool      `db:"allows_multiple_answers"`
	AllowsAnswerEdits     bool      `db:"allows_answer_edits"`
}
