package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

// ProfileRow represents a single PostgreSQL row containing the data of a profile
type ProfileRow struct {
	Address      string    `db:"address"`
	DTag         string    `db:"dtag"`
	Nickname     string    `db:"nickname"`
	Bio          string    `db:"bio"`
	ProfilePic   string    `db:"profile_pic"`
	CoverPic     string    `db:"cover_pic"`
	CreationTime time.Time `db:"creation_time"`
	Height       int64     `db:"height"`
}

// ConvertProfileRow converts the given row into a profile
func ConvertProfileRow(row ProfileRow) (*profilestypes.Profile, error) {
	address, err := sdk.AccAddressFromBech32(row.Address)
	if err != nil {
		return nil, err
	}

	profile, err := profilestypes.NewProfile(
		row.DTag,
		row.Nickname,
		row.Bio,
		profilestypes.NewPictures(row.ProfilePic, row.CoverPic),
		row.CreationTime,
		authtypes.NewBaseAccountWithAddress(address),
	)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
