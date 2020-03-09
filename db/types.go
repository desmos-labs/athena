package db

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/x/posts"
)

// PostRow represents a single PostgreSQL row containing the data of a Post
type PostRow struct {
	PostID         string    `db:"id"`
	ParentID       string    `db:"parent_id"`
	Message        string    `db:"message"`
	Created        time.Time `db:"created"`
	LastEdited     time.Time `db:"last_edited"`
	AllowsComments bool      `db:"allows_comments"`
	Subspace       string    `db:"subspace"`
	Creator        string    `db:"creator"`
	PollID         *uint64   `db:"poll_id"`
	OptionalData   string    `db:"optional_data"`
}

func ConvertPostRow(row PostRow) (*posts.Post, error) {

	// Parse the post id
	postID, err := posts.ParsePostID(row.PostID)
	if err != nil {
		return nil, err
	}

	// Parse the parent id
	parentID, err := posts.ParsePostID(row.ParentID)
	if err != nil {
		return nil, err
	}

	// Parse the creator
	creator, err := sdk.AccAddressFromBech32(row.Creator)
	if err != nil {
		return nil, err
	}

	// Parse the optional data
	var optionalData map[string]string
	err = json.Unmarshal([]byte(row.OptionalData), &optionalData)
	if err != nil {
		return nil, err
	}

	post := posts.NewPost(postID, parentID, row.Message, row.AllowsComments, row.Subspace, optionalData, row.Created, creator)
	post.LastEdited = row.LastEdited

	return &post, nil
}
