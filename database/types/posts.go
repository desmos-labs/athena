package types

import (
	"database/sql"
	"time"

	"github.com/desmos-labs/djuno/v2/types"

	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type PostRow struct {
	ParentID     sql.NullString `db:"parent_id"`
	Created      time.Time      `db:"created"`
	LastEdited   time.Time      `db:"last_edited"`
	PostID       string         `db:"id"`
	Message      string         `db:"message"`
	Subspace     string         `db:"subspace"`
	Creator      string         `db:"creator_address"`
	CommentState string         `db:"comments_state"`
	Hidden       bool           `db:"hidden"`
	Height       int64          `db:"height"`
}

// ConvertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func ConvertPostRow(
	row PostRow, attributes []poststypes.Attribute,
	attachments []poststypes.Attachment, poll *poststypes.Poll,
) (*types.Post, error) {
	var parentID string
	if row.ParentID.Valid {
		parentID = row.ParentID.String
	}

	state, err := poststypes.CommentsStateFromString(row.CommentState)
	if err != nil {
		return nil, err
	}

	return types.NewPost(
		poststypes.NewPost(
			row.PostID,
			parentID,
			row.Message,
			state,
			row.Subspace,
			attributes,
			attachments,
			poll,
			row.LastEdited,
			row.Created,
			row.Creator,
		),
		row.Height,
	), nil
}

// ---------------------------------------------------------------------------------------------------

// AttributeRow represents a single row inside the optional_data table
type AttributeRow struct {
	PostID string `db:"post_id"`
	Key    string `db:"key"`
	Value  string `db:"value"`
}

// ConvertAttributeRow converts the given rows into an OptionalData object
func ConvertAttributeRow(rows []AttributeRow) []poststypes.Attribute {
	var attributes = make([]poststypes.Attribute, len(rows))
	for index, row := range rows {
		attributes[index] = poststypes.NewAttribute(row.Key, row.Value)
	}
	return attributes
}

// ---------------------------------------------------------------------------------------------------

// AttachmentRow represents a single row inside the optional_data table
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
func ConvertPollRow(row PollRow, answers []poststypes.ProvidedAnswer) *poststypes.Poll {
	return poststypes.NewPoll(
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
func ConvertPollAnswerRows(rows []PollAnswerRow) []poststypes.ProvidedAnswer {
	answers := make([]poststypes.ProvidedAnswer, len(rows))
	for index, rows := range rows {
		answers[index] = poststypes.NewProvidedAnswer(rows.AnswerID, rows.Text)
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
