package database

import (
	"encoding/json"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"time"

	dbtypes "github.com/desmos-labs/djuno/database/types"
	"github.com/rs/zerolog/log"
)

// convertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func convertPostRow(postRow dbtypes.PostRow) (*poststypes.Post, error) {
	// Parse the optional data
	var optionalData map[string]string
	err := json.Unmarshal([]byte(postRow.OptionalData), &optionalData)
	if err != nil {
		return nil, err
	}

	post := poststypes.NewPost(
		postRow.PostID,
		postRow.ParentID,
		postRow.Message,
		postRow.AllowsComments,
		postRow.Subspace,
		optionalData,
		postRow.Created,
		postRow.Creator,
	)
	post.LastEdited = postRow.LastEdited

	return &post, nil
}

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post poststypes.Post) error {
	log.Info().Str("module", "posts").Str("post_id", post.PostID).Msg("saving post")

	err := db.savePostContent(post)
	if err != nil {
		return err
	}

	err = db.SavePollData(post.PostID, post.PollData)
	if err != nil {
		return err
	}

	// Save medias
	return db.saveAttachments(post.PostID, post.Attachments)
}

// savePostContent allows to store the content of the given post
func (db DesmosDb) savePostContent(post poststypes.Post) error {
	optionalDataBz, err := json.Marshal(post.OptionalData)
	if err != nil {
		return err
	}

	// Save the user
	err = db.SaveUserIfNotExisting(post.Creator)
	if err != nil {
		return err
	}

	// Save the post
	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator_address, optional_data, hidden)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err = db.Sql.Exec(
		postSqlStatement,
		post.PostID, post.ParentID, post.Message, post.Created, post.LastEdited, post.AllowsComments, post.Subspace,
		post.Creator, string(optionalDataBz), false,
	)
	return err
}

// saveAttachments allows to save the specified medias that are associated
// to the post having the given postID
func (db DesmosDb) saveAttachments(postID string, attachments []poststypes.Attachment) error {
	mediaQuery := `INSERT INTO media (post_id, uri, mime_type) VALUES ($1, $2, $3)`
	for _, media := range attachments {
		_, err := db.Sql.Exec(mediaQuery, postID, media.URI, media.MimeType)
		if err != nil {
			return err
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
	stmt = `DELETE FROM media WHERE post_id = $1`
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
	postSqlStatement := `SELECT * FROM post WHERE id = $1`

	var rows []dbtypes.PostRow
	err := db.sqlx.Select(&rows, postSqlStatement, id)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, nil
	}

	postRow := rows[0]
	return convertPostRow(postRow)
}
