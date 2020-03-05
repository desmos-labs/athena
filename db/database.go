package db

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/db/postgresql"
)

// DesmosDb represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type DesmosDb struct {
	*postgresql.Database
}

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post posts.Post) error {
	var pollID *uint64

	if post.PollData != nil {
		// Saving post's poll data before post to make possible the insertion of poll_id inside it
		statement := `
		INSERT INTO poll (question, end_date, open, allows_multiple_answers, allows_answer_edits)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
		`

		err := db.Sql.QueryRow(
			statement,
			post.PollData.Question, post.PollData.EndDate, post.PollData.Open, post.PollData.AllowsMultipleAnswers,
			post.PollData.AllowsAnswerEdits,
		).Scan(&pollID)

		if err != nil {
			return err
		}

		addPollAnswersSqlStatement := `
		INSERT INTO poll_answer(poll_id, answer_id, answer_text)
		VALUES($1, $2, $3)
		RETURNING id;
		`

		for _, answer := range post.PollData.ProvidedAnswers {
			err := db.Sql.QueryRow(
				addPollAnswersSqlStatement,
				pollID, answer.ID, answer.Text,
			).Scan()

			if err != nil {
				return err
			}
		}

	}

	// Saving Post
	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator, poll_id, optional_data)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;
    `

	// todo look how this is inserted in DB
	jsonB, _ := json.Marshal(post.OptionalData)

	err := db.Sql.QueryRow(
		postSqlStatement,
		post.PostID, post.ParentID, post.Message, post.Created, post.LastEdited, post.AllowsComments, post.Subspace,
		post.Creator.String(), pollID, string(jsonB),
	).Scan()
	if err != nil {
		return err
	}

	// Saving post's medias
	mediasSqlStatement := `
	INSERT INTO media (post_id, uri, mime_type)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	for _, media := range post.Medias {
		err = db.Sql.QueryRow(
			mediasSqlStatement,
			post.PostID, media.URI, media.MimeType,
		).Scan()
		if err != nil {
			return err
		}
	}

	return nil
}

// EditPost allows to properly edit the post having the given postID by setting the new
// given message and editDate
func (db DesmosDb) EditPost(postID posts.PostID, message string, editDate time.Time) error {
	statement := `
	UPDATE post 
	SET message = $1, last_edited = $2 
	WHERE id = $3
	RETURNING id;
	`

	return db.Sql.QueryRow(statement, message, editDate, postID).Scan()
}

// SaveReaction allows to save a new reaction for the given postID having the specified value and user
func (db DesmosDb) SaveReaction(postID posts.PostID, value string, user sdk.AccAddress) error {
	statement := `
	INSERT INTO reaction (post_id, owner, value)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	return db.Sql.QueryRow(statement, postID, user.String(), value).Scan()
}
