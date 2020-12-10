package types

import (
	"database/sql"
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type PostRow struct {
	PostID         string         `db:"id"`
	ParentID       sql.NullString `db:"parent_id"`
	Message        string         `db:"message"`
	Created        time.Time      `db:"created"`
	LastEdited     time.Time      `db:"last_edited"`
	AllowsComments bool           `db:"allows_comments"`
	Subspace       string         `db:"subspace"`
	Creator        string         `db:"creator_address"`
	PollID         *uint64        `db:"poll_id"`
	OptionalData   string         `db:"optional_data"`
	Hidden         bool           `db:"hidden"`
}

// ConvertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func ConvertPostRow(
	row PostRow, optionalData poststypes.OptionalData,
	attachments []poststypes.Attachment, poll *poststypes.PollData,
) poststypes.Post {
	var parentID string
	if row.ParentID.Valid {
		parentID = row.ParentID.String
	}

	return poststypes.NewPost(
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
	Uri      string `db:"uri"`
	MimeType string `db:"mime_type"`
}

// ConvertAttachmentRow converts the given row and tags into a proper Attachment object
func ConvertAttachmentRow(row AttachmentRow, tags []string) poststypes.Attachment {
	return poststypes.NewAttachment(row.Uri, row.MimeType, tags)
}

// AttachmentTagRow represents a single row of the attachment_tag table
type AttachmentTagRow struct {
	AttachmentID int    `db:"attachment_id"`
	Tag          string `db:"tag"`
}

// ConvertAttachmentTagRows converts the given rows into a slice of tags
func ConvertAttachmentTagRows(rows []AttachmentTagRow) []string {
	tags := make([]string, len(rows))
	for index, row := range rows {
		tags[index] = row.Tag
	}
	return tags
}
