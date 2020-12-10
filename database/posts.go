package database

import (
	"database/sql"
	"fmt"
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"
	"github.com/rs/zerolog/log"
)

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post poststypes.Post) error {
	log.Info().Str("module", "posts").Str("post_id", post.PostID).Msg("saving post")

	err := db.savePostContent(post)
	if err != nil {
		return err
	}

	err = db.saveOptionalData(post.PostID, post.OptionalData)
	if err != nil {
		return err
	}

	err = db.saveAttachments(post.PostID, post.Attachments)
	if err != nil {
		return err
	}

	err = db.SavePollData(post.PostID, post.PollData)
	if err != nil {
		return err
	}

	return err
}

// savePostContent allows to store the content of the given post
func (db DesmosDb) savePostContent(post poststypes.Post) error {
	// Save the user
	err := db.SaveUserIfNotExisting(post.Creator)
	if err != nil {
		return err
	}

	// Save the post
	stmt := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator_address, hidden)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// Convert the parent id string
	var parentID sql.NullString
	if len(post.ParentID) > 0 {
		parentID = sql.NullString{Valid: true, String: post.ParentID}
	}

	_, err = db.Sql.Exec(
		stmt,
		post.PostID, parentID, post.Message, post.Created, post.LastEdited, post.AllowsComments,
		post.Subspace, post.Creator, false,
	)
	return err
}

// saveOptionalData allows to save the specified optional data that are associated
// to the post having the given postID
func (db DesmosDb) saveOptionalData(postID string, data poststypes.OptionalData) error {
	stmt := `INSERT INTO optional_data (post_id, key, value) VALUES `
	var args []interface{}
	for index, entry := range data {
		oi := index * 3
		stmt += fmt.Sprintf("($%d, $%d, $%d),", oi+1, oi+2, oi+3)
		args = append(args, postID, entry.Key, entry.Value)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, args...)
	return err
}

// saveAttachments allows to save the specified medias that are associated
// to the post having the given postID
func (db DesmosDb) saveAttachments(postID string, attachments []poststypes.Attachment) error {
	for _, media := range attachments {
		// Insert the attachment
		var attachmentID int
		stmt := `INSERT INTO attachment (post_id, uri, mime_type) VALUES ($1, $2, $3) RETURNING id`
		err := db.Sqlx.QueryRow(stmt, postID, media.URI, media.MimeType).Scan(&attachmentID)
		if err != nil {
			return err
		}

		// Insert all the tags
		for _, tag := range media.Tags {
			err = db.SaveUserIfNotExisting(tag)
			if err != nil {
				return err
			}

			stmt = `INSERT INTO attachment_tag (attachment_id, tag) VALUES ($1, $2)`
			_, err = db.Sql.Exec(stmt, attachmentID, tag)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// EditPost allows to properly edit the post having the given postID by setting the new
// given message and editDate
func (db DesmosDb) EditPost(
	postID string, message string, attachments []poststypes.Attachment, poll *poststypes.PollData, editDate time.Time,
) error {
	stmt := `UPDATE post SET message = $1, last_edited = $2 WHERE id = $3`
	_, err := db.Sql.Exec(stmt, message, editDate, postID)
	if err != nil {
		return err
	}

	// Delete and store again the medias
	stmt = `DELETE FROM attachment WHERE post_id = $1`
	_, err = db.Sql.Exec(stmt, postID)
	if err != nil {
		return err
	}

	err = db.saveAttachments(postID, attachments)
	if err != nil {
		return err
	}

	// Delete and store again the poll data
	err = db.DeletePollData(postID)
	if err != nil {
		return err
	}

	err = db.SavePollData(postID, poll)
	return err
}

// GetPostByID returns the post having the specified id.
// If some error raised during the read, it is returned.
// If no post with the specified id is found, nil is returned instead.
func (db DesmosDb) GetPostByID(id string) (*poststypes.Post, error) {
	stmt := `SELECT * FROM post WHERE id = $1`

	var rows []dbtypes.PostRow
	err := db.Sqlx.Select(&rows, stmt, id)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, fmt.Errorf("no post with the given id found: %s", id)
	}

	row := rows[0]

	optionalData, err := db.getOptionalData(row.PostID)
	if err != nil {
		return nil, err
	}

	attachments, err := db.getAttachments(row.PostID)
	if err != nil {
		return nil, err
	}

	poll, err := db.GetPollByPostID(row.PostID)
	if err != nil {
		return nil, err
	}

	post := dbtypes.ConvertPostRow(row, optionalData, attachments, poll)
	return &post, nil
}

// getOptionalData returns all the optional data associated to the post having the given id
func (db DesmosDb) getOptionalData(postID string) (poststypes.OptionalData, error) {
	stmt := `SELECT * FROM optional_data WHERE post_id = $1`

	var rows []dbtypes.OptionalDataRow
	err := db.Sqlx.Select(&rows, stmt, postID)
	if err != nil {
		return nil, err
	}

	return dbtypes.ConvertOptionalDataRows(rows), nil
}

// getAttachments returns the attachments of the post having the given id
func (db DesmosDb) getAttachments(postID string) ([]poststypes.Attachment, error) {
	stmt := `SELECT * FROM attachment WHERE post_id = $1`

	var rows []dbtypes.AttachmentRow
	err := db.Sqlx.Select(&rows, stmt, postID)
	if err != nil {
		return nil, err
	}

	attachments := make([]poststypes.Attachment, len(rows))
	for i, row := range rows {
		var tagRows []dbtypes.AttachmentTagRow
		err := db.Sqlx.Select(&tagRows, `SELECT * FROM attachment_tag WHERE attachment_id  =$1`, row.ID)
		if err != nil {
			return nil, err
		}

		attachments[i] = dbtypes.ConvertAttachmentRow(row, dbtypes.ConvertAttachmentTagRows(tagRows))
	}

	return attachments, nil
}
