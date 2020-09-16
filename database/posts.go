package database

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"
	"github.com/rs/zerolog/log"
)

// convertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func convertPostRow(postRow dbtypes.PostRow, userRow *dbtypes.ProfileRow) (*poststypes.Post, error) {
	// Parse the parent id
	var err error
	var parentID poststypes.PostID
	if postRow.ParentID.Valid {
		parentID, err = poststypes.ParsePostID(postRow.ParentID.String)
		if err != nil {
			return nil, err
		}
	}

	// Parse the creator
	creator, err := sdk.AccAddressFromBech32(userRow.Address)
	if err != nil {
		return nil, err
	}

	// Parse the optional data
	var optionalData map[string]string
	err = json.Unmarshal([]byte(postRow.OptionalData), &optionalData)
	if err != nil {
		return nil, err
	}

	post := poststypes.NewPost(
		parentID,
		postRow.Message,
		postRow.AllowsComments,
		postRow.Subspace,
		optionalData,
		postRow.Created,
		creator,
	)
	post.LastEdited = postRow.LastEdited

	return &post, nil
}

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post poststypes.Post) error {
	log.Info().Str("module", "poststypes").
		Str("post_id", post.PostID.String()).
		Msg("saving post")

	err := db.savePostContent(post)
	if err != nil {
		return err
	}

	err = db.SavePollData(post.PostID, post.PollData)
	if err != nil {
		return err
	}

	// Save medias
	return db.saveMedias(post.PostID, post.Attachments)
}

// savePostContent allows to store the content of the given post
func (db DesmosDb) savePostContent(post poststypes.Post) error {
	optionalDataBz, err := json.Marshal(post.OptionalData)
	if err != nil {
		return err
	}

	// Save the user
	if err := db.SaveUserIfNotExisting(post.Creator); err != nil {
		return err
	}

	// Save the post
	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator_address, optional_data, hidden)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	ON CONFLICT (id) DO NOTHING`

	var parentId *string
	if post.ParentID.Valid() {
		parentIdString := post.ParentID.String()
		parentId = &parentIdString
	}

	_, err = db.Sql.Exec(
		postSqlStatement,
		post.PostID.String(), parentId, post.Message, post.Created, post.LastEdited, post.AllowsComments, post.Subspace,
		post.Creator.String(), string(optionalDataBz), false,
	)
	return err
}

// saveMedias allows to save the specified medias that are associated
// to the post having the given postID
func (db DesmosDb) saveMedias(postID poststypes.PostID, medias poststypes.Attachments) error {
	stmt := `INSERT INTO media (post_id, uri, mime_type) 
			 VALUES ($1, $2, $3) 
			 ON CONFLICT ON CONSTRAINT unique_post_media DO NOTHING`

	for _, media := range medias {
		_, err := db.Sql.Exec(stmt, postID.String(), media.URI, media.MimeType)
		if err != nil {
			return err
		}
	}

	return nil
}

// EditPost allows to properly edit the post having the given postID by setting the new
// given message and editDate
func (db DesmosDb) EditPost(postID poststypes.PostID, message string, editDate time.Time) error {
	statement := `UPDATE post SET message = $1, last_edited = $2 WHERE id = $3`
	_, err := db.Sql.Exec(statement, message, editDate, postID)
	return err
}

// GetPostByID returns the post having the specified id.
// If some error raised during the read, it is returned.
// If no post with the specified id is found, nil is returned instead.
func (db DesmosDb) GetPostByID(id poststypes.PostID) (*poststypes.Post, error) {
	postSqlStatement := `SELECT * FROM post WHERE id = $1`

	var rows []dbtypes.PostRow
	err := db.Sqlx.Select(&rows, postSqlStatement, id)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, nil
	}

	postRow := rows[0]

	// Find the user
	addr, err := sdk.AccAddressFromBech32(postRow.Creator)
	if err != nil {
		return nil, err
	}

	userRow, err := db.GetUserByAddress(addr)
	if err != nil {
		return nil, err
	}

	return convertPostRow(postRow, userRow)
}
