package db

import (
	"database/sql"
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/rs/zerolog/log"
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
	CreatorID      *uint64        `db:"creator_id"`
	PollID         *uint64        `db:"poll_id"`
	OptionalData   string         `db:"optional_data"`
	Hidden         bool           `db:"hidden"`
}

// ConvertPostRow takes the given postRow and userRow and merges the data contained inside them to create a Post.
func ConvertPostRow(postRow PostRow, userRow *UserRow) (*posts.Post, error) {

	// Parse the post id
	postID, err := posts.ParsePostID(postRow.PostID)
	if err != nil {
		return nil, err
	}

	// Parse the parent id

	var parentID posts.PostID
	if postRow.ParentID.Valid {
		parentID, err = posts.ParsePostID(postRow.ParentID.String)
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

	post := posts.NewPost(postID, parentID, postRow.Message, postRow.AllowsComments, postRow.Subspace,
		optionalData, postRow.Created, creator)
	post.LastEdited = postRow.LastEdited

	return &post, nil
}

// GetPostByID returns the post having the specified id.
// If some error raised during the read, it is returned.
// If no post with the specified id is found, nil is returned instead.
func (db DesmosDb) GetPostByID(id posts.PostID) (*posts.Post, error) {
	postSqlStatement := `SELECT * FROM post WHERE id = $1`

	var rows []PostRow
	err := db.sqlx.Select(&rows, postSqlStatement, id)
	if err != nil {
		return nil, err
	}

	// No post found
	if len(rows) == 0 {
		return nil, nil
	}

	postRow := rows[0]

	// Find the user
	userRow, err := db.GetUserById(postRow.CreatorID)
	if err != nil {
		return nil, err
	}

	return ConvertPostRow(postRow, userRow)
}

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post posts.Post) error {
	log.Info().Str("post_id", post.PostID.String()).Msg("saving post")

	user, err := db.SaveUserIfNotExisting(post.Creator)
	if err != nil {
		return err
	}

	pollID, err := db.SavePollData(post.PollData)
	if err != nil {
		return err
	}

	err = db.savePostContent(post, user.Id, pollID)
	if err != nil {
		return err
	}

	err = db.saveComment(post)
	if err != nil {
		return err
	}

	// Save medias
	return db.saveMedias(post.PostID, post.Medias)
}

// savePostContent allows to store the content of the given post, which has been created
// from the user having the specified id and contains a poll with the given id
func (db DesmosDb) savePostContent(post posts.Post, userID *uint64, pollID *uint64) error {
	optionalDataBz, err := json.Marshal(post.OptionalData)
	if err != nil {
		return err
	}

	// Saving Post
	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, 
	                  creator_id, poll_id, optional_data, hidden)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	var parentId *string
	// TODO: Remove the second part of the check once the invariants have been implemented
	if post.ParentID.Valid() && post.ParentID < post.PostID {
		parentIdString := post.ParentID.String()
		parentId = &parentIdString
	}

	_, err = db.Sql.Exec(
		postSqlStatement,
		post.PostID.String(), parentId, post.Message, post.Created, post.LastEdited, post.AllowsComments, post.Subspace,
		userID, pollID, string(optionalDataBz), false,
	)
	return err
}

func (db DesmosDb) saveComment(post posts.Post) error {
	// TODO: Remove the second part of the check once the invariants have been implemented
	if !post.ParentID.Valid() || post.ParentID > post.PostID {
		return nil
	}

	commentSqlStatement := `INSERT INTO comment (post_id, comment_id) VALUES ($1, $2)`
	_, err := db.Sql.Exec(commentSqlStatement, post.ParentID.String(), post.PostID.String())
	return err
}

// saveMedias allows to save the specified medias that are associated
// to the post having the given postID
func (db DesmosDb) saveMedias(postID posts.PostID, medias posts.PostMedias) error {
	mediaQuery := `INSERT INTO media (post_id, uri, mime_type) VALUES ($1, $2, $3)`
	for _, media := range medias {
		_, err := db.Sql.Exec(mediaQuery, postID.String(), media.URI, media.MimeType)
		if err != nil {
			return err
		}
	}

	return nil
}

// EditPost allows to properly edit the post having the given postID by setting the new
// given message and editDate
func (db DesmosDb) EditPost(postID posts.PostID, message string, editDate time.Time) error {
	statement := `UPDATE post SET message = $1, last_edited = $2 WHERE id = $3`
	_, err := db.Sql.Exec(statement, message, editDate, postID)
	return err
}
