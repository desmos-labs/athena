package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	poststypes "github.com/desmos-labs/desmos/v6/x/posts/types"

	dbtypes "github.com/desmos-labs/athena/database/types"
	"github.com/desmos-labs/athena/types"
)

// getPostRowID returns the row_id of the post having the given data
func (db *Db) getPostRowID(subspaceID uint64, postID uint64) (sql.NullInt64, error) {
	stmt := `SELECT row_id FROM post WHERE subspace_id = $1 and id = $2`

	var rowID int64
	err := db.SQL.QueryRow(stmt, subspaceID, postID).Scan(&rowID)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.NullInt64{Int64: 0, Valid: false}, nil
	}

	return sql.NullInt64{Int64: rowID, Valid: true}, err
}

// SavePost stores the given post inside the database
func (db *Db) SavePost(post types.Post) error {
	// Get the section row id
	sectionRowID, err := db.getSectionRowID(post.SubspaceID, post.SectionID)
	if err != nil {
		return err
	}

	// Get the conversation row id
	conversationRowID, err := db.getPostRowID(post.SubspaceID, post.ConversationID)
	if err != nil {
		return err
	}

	// Insert the post
	stmt := `
INSERT INTO post (subspace_id, section_row_id, id, external_id, text, author_address, owner_address, conversation_row_id, reply_settings, creation_date, last_edited_date, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT ON CONSTRAINT unique_subspace_post DO UPDATE 
    SET subspace_id = excluded.subspace_id,
        section_row_id = excluded.section_row_id,
        id = excluded.id,
        external_id = excluded.external_id,
        text = excluded.text,
        author_address = excluded.author_address,
        owner_address = excluded.owner_address,
        conversation_row_id = excluded.conversation_row_id,
        reply_settings = excluded.reply_settings,
        creation_date = excluded.creation_date,
        last_edited_date = excluded.last_edited_date,
        height = excluded.height
WHERE post.height <= excluded.height
RETURNING row_id`

	var rowID uint64
	err = db.SQL.QueryRow(stmt,
		post.SubspaceID,
		sectionRowID,
		post.ID,
		dbtypes.ToNullString(post.ExternalID),
		dbtypes.ToNullString(post.Text),
		post.Author,
		dbtypes.ToNullString(post.Owner),
		conversationRowID,
		post.ReplySettings.String(),
		post.CreationDate,
		dbtypes.ToNullTime(post.LastEditedDate),
		post.Height,
	).Scan(&rowID)
	if err != nil {
		return err
	}

	// Insert the entities
	err = db.savePostEntities(rowID, post.Entities)
	if err != nil {
		return err
	}

	// Insert the tags
	err = db.savePostTags(rowID, post.Tags)
	if err != nil {
		return err
	}

	// Insert the reference
	err = db.savePostReferences(post.SubspaceID, rowID, post.ReferencedPosts)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) savePostEntities(postRowID uint64, entities *poststypes.Entities) error {
	if entities == nil {
		return nil
	}

	err := db.savePostHashtags(postRowID, entities.Hashtags)
	if err != nil {
		return err
	}

	err = db.savePostMentions(postRowID, entities.Mentions)
	if err != nil {
		return err
	}

	err = db.savePostURLs(postRowID, entities.Urls)
	if err != nil {
		return err
	}

	return nil
}

func (db *Db) savePostHashtags(postRowID uint64, hashtags []poststypes.TextTag) error {
	// Delete all hashtags first
	stmt := `DELETE FROM post_hashtag WHERE post_row_id = $1`
	_, err := db.SQL.Exec(stmt, postRowID)
	if err != nil {
		return err
	}

	if hashtags == nil {
		return nil
	}

	stmt = `INSERT INTO post_hashtag (post_row_id, start_index, end_index, tag) VALUES `

	var vars []interface{}
	for i, hashtag := range hashtags {
		ei := i * 4
		stmt += fmt.Sprintf(`($%d, $%d, $%d, $%d),`, ei+1, ei+2, ei+3, ei+4)
		vars = append(vars, postRowID, hashtag.Start, hashtag.End, hashtag.Tag)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, vars...)
	return err
}

func (db *Db) savePostMentions(postRowID uint64, mentions []poststypes.TextTag) error {
	// Delete all mentions first
	stmt := `DELETE FROM post_mention WHERE post_row_id = $1`
	_, err := db.SQL.Exec(stmt, postRowID)
	if err != nil {
		return err
	}

	if mentions == nil {
		return nil
	}

	stmt = `INSERT INTO post_mention (post_row_id, start_index, end_index, mention_address) VALUES `

	var vars []interface{}
	for i, mention := range mentions {
		ei := i * 4
		stmt += fmt.Sprintf(`($%d, $%d, $%d, $%d),`, ei+1, ei+2, ei+3, ei+4)
		vars = append(vars, postRowID, mention.Start, mention.End, mention.Tag)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, vars...)
	return err
}

func (db *Db) savePostURLs(postRowID uint64, urls []poststypes.Url) error {
	// Delete all urls first
	stmt := `DELETE FROM post_url WHERE post_row_id = $1`
	_, err := db.SQL.Exec(stmt, postRowID)
	if err != nil {
		return err
	}

	if urls == nil {
		return nil
	}

	// Save the urls
	stmt = `INSERT INTO post_url (post_row_id, start_index, end_index, url, display_value) VALUES `

	var vars []interface{}
	for i, url := range urls {
		ei := i * 5
		stmt += fmt.Sprintf(`($%d, $%d, $%d, $%d, $%d),`, ei+1, ei+2, ei+3, ei+4, ei+5)
		vars = append(vars, postRowID, url.Start, url.End, url.Url, url.DisplayUrl)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, vars...)
	return err
}

func (db *Db) savePostTags(postRowID uint64, tags []string) error {
	// Delete all tags first
	stmt := `DELETE FROM post_tag WHERE post_row_id = $1`
	_, err := db.SQL.Exec(stmt, postRowID)
	if err != nil {
		return err
	}

	if tags == nil {
		return nil
	}

	// Save the urls
	stmt = `INSERT INTO post_tag (post_row_id, tag) VALUES `

	var vars []interface{}
	for i, tag := range tags {
		ei := i * 2
		stmt += fmt.Sprintf(`($%d, $%d),`, ei+1, ei+2)
		vars = append(vars, postRowID, tag)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, vars...)
	return err
}

func (db *Db) savePostReferences(subspaceID uint64, postRowID uint64, references []poststypes.PostReference) error {
	// Delete all references first
	stmt := `DELETE FROM post_reference WHERE post_row_id = $1`
	_, err := db.SQL.Exec(stmt, postRowID)
	if err != nil {
		return err
	}

	if len(references) == 0 {
		return nil
	}

	stmt = `INSERT INTO post_reference (post_row_id, type, reference_row_id, position_index) VALUES `

	var vars []interface{}
	for i, ref := range references {
		referenceRowID, err := db.getPostRowID(subspaceID, ref.PostID)
		if err != nil {
			return err
		}

		ei := i * 4
		stmt += fmt.Sprintf(`($%d, $%d, $%d, $%d),`, ei+1, ei+2, ei+3, ei+4)
		vars = append(vars, postRowID, ref.Type.String(), referenceRowID, ref.Position)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, vars...)
	return err
}

// HasPost returns true if the post with the given id exists inside the database
func (db *Db) HasPost(height int64, subspaceID uint64, postID uint64) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM post WHERE subspace_id = $1 AND id = $2 AND height <= $3)`
	var exists bool
	err := db.SQL.QueryRow(stmt, subspaceID, postID, height).Scan(&exists)
	return exists, err
}

// DeletePost removes the post with the given details from the database
func (db *Db) DeletePost(height int64, subspaceID uint64, postID uint64) error {
	stmt := `DELETE FROM post WHERE subspace_id = $1 AND id = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, subspaceID, postID, height)
	return err
}

// DeleteAllPosts removes all the posts for the given subspace from the database
func (db *Db) DeleteAllPosts(height int64, subspaceID uint64) error {
	stmt := `DELETE FROM post WHERE height <= $1 AND subspace_id = $2`
	_, err := db.SQL.Exec(stmt, height, subspaceID)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SavePostTx stores the given transaction into the database
func (db *Db) SavePostTx(tx types.PostTransaction) error {
	postRowID, err := db.getPostRowID(tx.SubspaceID, tx.PostID)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO post_transaction (post_row_id, hash) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err = db.SQL.Exec(stmt, postRowID, tx.Hash)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

func (db *Db) getAttachmentRowID(subspaceID uint64, postID uint64, attachmentID uint32) (int64, error) {
	stmt := `
SELECT row_id FROM post_attachment WHERE post_row_id = (
	SELECT row_id FROM post WHERE subspace_id = $1 AND id = $2
) and id = $3`

	var rowID int64
	err := db.SQL.QueryRow(stmt, subspaceID, postID, attachmentID).Scan(&rowID)
	if errors.Is(err, sql.ErrNoRows) {
		return rowID, nil
	}

	return rowID, err
}

// SavePostAttachment stores the given attachment inside the database
func (db *Db) SavePostAttachment(attachment types.PostAttachment) error {
	postRowID, err := db.getPostRowID(attachment.SubspaceID, attachment.PostID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO post_attachment (post_row_id, id, content, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_post_attachment DO UPDATE 
    SET content = excluded.content,
        height = excluded.height
WHERE post_attachment.height <= excluded.height`

	contentBz, err := db.cdc.MarshalJSON(attachment.Content)
	if err != nil {
		return fmt.Errorf("failed to json encode attachment content: %s", err)
	}

	_, err = db.SQL.Exec(stmt,
		postRowID,
		attachment.ID,
		string(contentBz),
		attachment.Height,
	)
	return err
}

// DeletePostAttachment removes the given post attachment from the database
func (db *Db) DeletePostAttachment(height int64, subspaceID uint64, postID uint64, attachmentID uint32) error {
	stmt := `
DELETE FROM post_attachment WHERE post_row_id = (
	SELECT row_id FROM post WHERE subspace_id = $1 AND id = $2
) AND id = $3 AND height <= $4`
	_, err := db.SQL.Exec(stmt, subspaceID, postID, attachmentID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SavePollAnswer stores the given answer inside the database
func (db *Db) SavePollAnswer(answer types.PollAnswer) error {
	attachmentRowID, err := db.getAttachmentRowID(answer.SubspaceID, answer.PostID, answer.PollID)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO poll_answer (attachment_row_id, answers_indexes, user_address, height)
VALUES ($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT unique_user_answer DO UPDATE 
    SET answers_indexes = excluded.answers_indexes,
        user_address = excluded.user_address,
        height = excluded.height
WHERE poll_answer.height <= excluded.height`
	_, err = db.SQL.Exec(stmt,
		attachmentRowID,
		answer.AnswersIndexes,
		answer.User,
		answer.Height,
	)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SavePostsParams stores the given params inside the database
func (db *Db) SavePostsParams(params types.PostsParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling reports params: %s", err)
	}

	stmt := `
INSERT INTO posts_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE posts_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing reports params: %s", err)
	}

	return nil
}
