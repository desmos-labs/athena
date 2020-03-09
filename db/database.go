package db

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DesmosDb represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type DesmosDb struct {
	postgresql.Database
	sqlx *sqlx.DB
}

// Builder allows to create a new DesmosDb instance implementing the database.Builder type
func Builder(cfg config.Config, codec *codec.Codec) (*db.Database, error) {
	database, err := postgresql.Builder(cfg, codec)
	if err != nil {
		return nil, err
	}

	psqlDb, _ := (*database).(postgresql.Database)
	var desmosDb db.Database = DesmosDb{
		Database: psqlDb,
		sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}

	return &desmosDb, nil
}

// SavePost allows to store the given post inside the database properly.
func (db DesmosDb) SavePost(post posts.Post) error {
	var pollID *uint64

	// TODO: Split this
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

		pollQuery := `INSERT INTO poll_answer(poll_id, answer_id, answer_text) VALUES($1, $2, $3)`
		for _, answer := range post.PollData.ProvidedAnswers {
			_, err := db.Sql.Exec(pollQuery, pollID, answer.ID, answer.Text)
			if err != nil {
				return err
			}
		}

	}

	// Saving Post
	postSqlStatement := `
	INSERT INTO post (id, parent_id, message, created, last_edited, allows_comments, subspace, creator, poll_id, optional_data)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	jsonB, err := json.Marshal(post.OptionalData)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(
		postSqlStatement,
		post.PostID, post.ParentID, post.Message, post.Created, post.LastEdited, post.AllowsComments, post.Subspace,
		post.Creator.String(), pollID, string(jsonB),
	)
	if err != nil {
		return err
	}

	// Saving post's medias
	// TODO: Split this
	mediaQuery := `INSERT INTO media (post_id, uri, mime_type) VALUES ($1, $2, $3)`
	for _, media := range post.Medias {
		_, err = db.Sql.Exec(mediaQuery, post.PostID, media.URI, media.MimeType)
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

	return ConvertPostRow(rows[0])
}

// SaveReaction allows to save a new reaction for the given postID having the specified value and user
func (db DesmosDb) SaveReaction(postID posts.PostID, reaction posts.Reaction) (*posts.Reaction, error) {
	statement := `INSERT INTO reaction (post_id, owner, value) VALUES ($1, $2, $3)`
	_, err := db.Sql.Exec(statement, postID, reaction.Owner.String(), reaction.Value)
	return &reaction, err
}

// RemoveReaction allows to remove an already existing reaction for the post having the given postID,
// the given reaction and from the specified user.
func (db DesmosDb) RemoveReaction(postID posts.PostID, reaction string, user sdk.AccAddress) error {
	statement := `DELETE FROM reaction WHERE post_id = $1 AND owner = $2 AND reaction = $3;`
	_, err := db.Sql.Exec(statement, postID, user.String(), reaction)
	return err
}

// SavePollAnswer allows to save the given answers from the specified user for the poll
// post having the specified postID.
func (db DesmosDb) SavePollAnswer(postID posts.PostID, answer posts.UserAnswer) error {
	statement := `INSERT INTO user_poll_answer (poll_id, answers, user_address) VALUES ($1, $2, $3)`
	_, err := db.Sql.Exec(statement, postID, pq.Array(answer.Answers), answer.User.String())
	return err
}
