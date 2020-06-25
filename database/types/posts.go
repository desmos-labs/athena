package types

import (
	"database/sql"
	"time"
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
	Creator        string         `db:"creator"`
	PollID         *uint64        `db:"poll_id"`
	OptionalData   string         `db:"optional_data"`
	Hidden         bool           `db:"hidden"`
}
