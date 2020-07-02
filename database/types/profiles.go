package types

import (
	"database/sql"
)

// ProfileRow represents a single PostgreSQL row containing the data of a profile
type ProfileRow struct {
	Address      string         `db:"address"`
	DTag         sql.NullString `db:"dtag"`
	Moniker      sql.NullString `db:"moniker"`
	Bio          sql.NullString `db:"bio"`
	ProfilePic   sql.NullString `db:"profile_pic"`
	CoverPic     sql.NullString `db:"cover_pic"`
	CreationDate sql.NullTime   `db:"creation_date"`
}
