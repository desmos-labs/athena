package types

import (
	"database/sql"
	"time"

	"github.com/desmos-labs/djuno/types"

	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type PostRow struct {
	ParentID       sql.NullString `db:"parent_id"`
	Created        time.Time      `db:"created"`
	LastEdited     time.Time      `db:"last_edited"`
	PostID         string         `db:"id"`
	Message        string         `db:"message"`
	Subspace       string         `db:"subspace"`
	Creator        string         `db:"creator_address"`
	AllowsComments bool           `db:"allows_comments"`
	Hidden         bool           `db:"hidden"`
	Height         int64          `db:"height"`
}

// ConvertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func ConvertPostRow(
	row PostRow, optionalData poststypes.OptionalData,
	attachments []poststypes.Attachment, poll *poststypes.PollData,
) *types.Post {
	var parentID string
	if row.ParentID.Valid {
		parentID = row.ParentID.String
	}

	return types.NewPost(
		poststypes.NewPost(
			row.PostID,
			parentID,
			row.Message,
			row.AllowsComments,
			row.Subspace,
			optionalData,
			attachments,
			poll,
			row.LastEdited,
			row.Created,
			row.Creator,
		),
		row.Height,
	)
}

// ---------------------------------------------------------------------------------------------------

// OptionalDataRow represents a single row inside the optional_data table
type OptionalDataRow struct {
	PostID string `db:"post_id"`
	Key    string `db:"key"`
	Value  string `db:"value"`
}

// ConvertOptionalDataRows converts the given rows into an OptionalData object
func ConvertOptionalDataRows(rows []OptionalDataRow) poststypes.OptionalData {
	attachments := make(poststypes.OptionalData, len(rows))
	for index, row := range rows {
		attachments[index] = poststypes.NewOptionalDataEntry(row.Key, row.Value)
	}
	return attachments
}

// ---------------------------------------------------------------------------------------------------

// OptionalDataRow represents a single row inside the optional_data table
type AttachmentRow struct {
	ID       int    `db:"id"`
	PostID   string `db:"post_id"`
	URI      string `db:"uri"`
	MimeType string `db:"mime_type"`
}

// ConvertAttachmentRow converts the given row and tags into a proper Attachment object
func ConvertAttachmentRow(row AttachmentRow, tags []string) poststypes.Attachment {
	return poststypes.NewAttachment(row.URI, row.MimeType, tags)
}

// AttachmentTagRow represents a single row of the attachment_tag table
type AttachmentTagRow struct {
	AttachmentID int    `db:"attachment_id"`
	Tag          string `db:"tag_address"`
}

// ConvertAttachmentTagRows converts the given rows into a slice of tags
func ConvertAttachmentTagRows(rows []AttachmentTagRow) []string {
	tags := make([]string, len(rows))
	for index, row := range rows {
		tags[index] = row.Tag
	}
	return tags
}

// ---------------------------------------------------------------------------------------------------

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
	Height          int64  `db:"height"`
}

func NewUserPollAnswerRow(pollID int, answer string, answerer string, height int64) UserPollAnswerRow {
	return UserPollAnswerRow{
		PollID:          pollID,
		Answer:          answer,
		AnswererAddress: answerer,
		Height:          height,
	}
}

// Equal tells whether r and s contain the same data
func (r UserPollAnswerRow) Equal(s UserPollAnswerRow) bool {
	return r.PollID == s.PollID &&
		r.Answer == s.Answer &&
		r.AnswererAddress == s.AnswererAddress &&
		r.Height == s.Height
}
