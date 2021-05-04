package database

import (
	"database/sql"
	"fmt"

	"github.com/desmos-labs/djuno/x/posts/types"

	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"

	"github.com/rs/zerolog/log"

	dbtypes "github.com/desmos-labs/djuno/database/types"
)

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post types.Post) error {
	log.Info().Str("module", "posts").Str("post_id", post.PostId).Msg("saving post")

	// Delete any previous posts
	stmt := `DELETE FROM post WHERE id = $1 AND height <= $2`
	_, err := db.Sql.Exec(stmt, post.PostId, post.Height)
	if err != nil {
		return err
	}

	err = db.savePostContent(post)
	if err != nil {
		return err
	}

	err = db.saveOptionalData(post.PostId, post.OptionalData)
	if err != nil {
		return err
	}

	err = db.saveAttachments(post.Height, post.PostId, post.Attachments)
	if err != nil {
		return err
	}

	err = db.savePollData(post.PostId, post.PollData)
	if err != nil {
		return err
	}

	return err
}

// savePostContent allows to store the content of the given post
func (db DesmosDb) savePostContent(post types.Post) error {
	// Save the user
	err := db.SaveUserIfNotExisting(post.Creator, post.Height)
	if err != nil {
		return err
	}

	// Save the post
	stmt := `
INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator_address, hidden, height)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// Convert the parent id string
	var parentID sql.NullString
	if len(post.ParentId) > 0 {
		parentID = sql.NullString{Valid: true, String: post.ParentId}
	}

	_, err = db.Sql.Exec(
		stmt,
		post.PostId, parentID, post.Message, post.Created, post.LastEdited, post.AllowsComments,
		post.Subspace, post.Creator, false, post.Height,
	)
	return err
}

// saveOptionalData allows to save the specified optional data that are associated
// to the post having the given postID
func (db DesmosDb) saveOptionalData(postID string, data poststypes.OptionalData) error {
	stmt := `INSERT INTO post_optional_data_entry (post_id, key, value) VALUES `
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
func (db DesmosDb) saveAttachments(height int64, postID string, attachments []poststypes.Attachment) error {
	for _, media := range attachments {
		// Insert the attachment
		var attachmentID int
		stmt := `INSERT INTO post_attachment (post_id, uri, mime_type) VALUES ($1, $2, $3) RETURNING id`
		err := db.Sqlx.QueryRow(stmt, postID, media.URI, media.MimeType).Scan(&attachmentID)
		if err != nil {
			return err
		}

		// Insert all the tags
		for _, tag := range media.Tags {
			err = db.SaveUserIfNotExisting(tag, height)
			if err != nil {
				return err
			}

			stmt = `INSERT INTO post_attachment_tag (attachment_id, tag_address) VALUES ($1, $2)`
			_, err = db.Sql.Exec(stmt, attachmentID, tag)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// savePollData allows to properly store the given poll inside the database, returning the
// id of the newly created (or updated) row inside the database itself.
// If the given poll is nil, it will not be inserted and nil will be returned as the id.
func (db DesmosDb) savePollData(postID string, poll *poststypes.PollData) error {
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

// ---------------------------------------------------------------------------------------------------

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

	poll, err := db.getPollData(row.PostID)
	if err != nil {
		return nil, err
	}

	post := dbtypes.ConvertPostRow(row, optionalData, attachments, poll)
	return &post, nil
}

// getOptionalData returns all the optional data associated to the post having the given id
func (db DesmosDb) getOptionalData(postID string) (poststypes.OptionalData, error) {
	stmt := `SELECT * FROM post_optional_data_entry WHERE post_id = $1`

	var rows []dbtypes.OptionalDataRow
	err := db.Sqlx.Select(&rows, stmt, postID)
	if err != nil {
		return nil, err
	}

	return dbtypes.ConvertOptionalDataRows(rows), nil
}

// getAttachments returns the attachments of the post having the given id
func (db DesmosDb) getAttachments(postID string) ([]poststypes.Attachment, error) {
	stmt := `SELECT * FROM post_attachment WHERE post_id = $1`

	var rows []dbtypes.AttachmentRow
	err := db.Sqlx.Select(&rows, stmt, postID)
	if err != nil {
		return nil, err
	}

	attachments := make([]poststypes.Attachment, len(rows))
	for i, row := range rows {
		var tagRows []dbtypes.AttachmentTagRow
		err := db.Sqlx.Select(&tagRows, `SELECT * FROM post_attachment_tag WHERE attachment_id  =$1`, row.ID)
		if err != nil {
			return nil, err
		}

		attachments[i] = dbtypes.ConvertAttachmentRow(row, dbtypes.ConvertAttachmentTagRows(tagRows))
	}

	return attachments, nil
}

// getPollData returns the poll row associated to the post having the specified id.
// If the post with the same id has no poll associated to it, nil is returned instead.
func (db DesmosDb) getPollData(postID string) (*poststypes.PollData, error) {
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

// ---------------------------------------------------------------------------------------------------

// SaveUserPollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SaveUserPollAnswer(answer types.UserPollAnswer) error {
	err := db.SaveUserIfNotExisting(answer.User, answer.Height)
	if err != nil {
		return err
	}

	statement := `
INSERT INTO user_poll_answer (poll_id, answer, answerer_address, height) 
VALUES ((SELECT id FROM poll WHERE post_id = $1), $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_user_answer DO UPDATE 
    SET answer = excluded.answer, 
    	answerer_address = excluded.answerer_address,
    	height = excluded.height
WHERE user_poll_answer.height <= excluded.height`

	for _, answerText := range answer.Answers {
		_, err = db.Sql.Exec(statement, answer.PostID, answerText, answer.User, answer.Height)
		if err != nil {
			return err
		}
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------

// SaveReport allows to store the given report properly
func (db DesmosDb) SaveReport(report types.Report) error {
	err := db.SaveUserIfNotExisting(report.User, report.Height)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO report(post_id, type, message, reporter_address, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT ON CONSTRAINT unique_report DO UPDATE 
    SET post_id = excluded.post_id, 
        type = excluded.type,
        message = excluded.message,
        reporter_address = excluded.reporter_address, 
        height = excluded.height
WHERE report.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, report.PostId, report.Type, report.Message, report.User, report.Height)
	return err
}
