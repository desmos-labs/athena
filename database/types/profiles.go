package types

import (
	"time"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

// ProfileRow represents a single PostgreSQL row containing the data of a profile
type ProfileRow struct {
	Address      string    `db:"address"`
	DTag         string    `db:"dtag"`
	Moniker      string    `db:"moniker"`
	Bio          string    `db:"bio"`
	ProfilePic   string    `db:"profile_pic"`
	CoverPic     string    `db:"cover_pic"`
	CreationTime time.Time `db:"creation_time"`
}

// ConvertProfileRow converts the given row into a profile
func ConvertProfileRow(row ProfileRow) profilestypes.Profile {
	return profilestypes.NewProfile(
		row.DTag,
		row.Moniker,
		row.Bio,
		profilestypes.NewPictures(row.ProfilePic, row.CoverPic),
		row.CreationTime,
		row.Address,
	)
}
