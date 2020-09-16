package database_test

import (
	dbtypes "github.com/desmos-labs/djuno/database/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

func newStrPtr(value string) *string {
	return &value
}

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	addr, err := sdk.AccAddressFromBech32("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	err = suite.database.SaveUserIfNotExisting(addr)
	suite.Require().NoError(err, "storing of address should return no error")

	err = suite.database.SaveUserIfNotExisting(addr)
	suite.Require().NoError(err, "storing address second time should return no error")

	var rows []dbtypes.ProfileRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM profile")
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	suite.Require().Equal(addr.String(), rows[0].Address)
}

func (suite *DbTestSuite) TestDesmosDB_UpdateProfileDTag() {
	user, err := sdk.AccAddressFromBech32("cosmos15vqfryzfgaznd45cef25jlpn2wqdjs5a53ws45")
	suite.Require().NoError(err)

	err = suite.database.UpdateProfileDTag(user, "moniker")
	suite.Require().NoError(err)

	var rows []dbtypes.ProfileRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM profile WHERE address = $1", user.String())
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].DTag.Valid)
	suite.Require().Equal("moniker", rows[0].DTag.String)

	err = suite.database.UpdateProfileDTag(user, "new-dtag")
	suite.Require().NoError(err)

	var newRows []dbtypes.ProfileRow
	err = suite.database.Sqlx.Select(&newRows, "SELECT * FROM profile WHERE address = $1", user.String())
	suite.Require().NoError(err)

	suite.Require().Len(newRows, 1)
	suite.Require().True(newRows[0].DTag.Valid)
	suite.Require().Equal("new-dtag", newRows[0].DTag.String)
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	creator, err := sdk.AccAddressFromBech32("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)

	creationDate, err := time.Parse(time.RFC3339, "2020-01-01T12:00:00Z")
	suite.Require().NoError(err)

	profile := profilestypes.NewProfile("dtag", creator, creationDate).
		WithMoniker(newStrPtr("profile-moniker"))

	// Save data
	err = suite.database.SaveProfile(profile)
	suite.Require().NoError(err)

	stored, err := suite.database.GetUserByAddress(profile.Creator)
	suite.Require().NoError(err)
	suite.Require().Equal(profile.DTag, stored.DTag.String)
	suite.Require().Equal(*profile.Moniker, stored.Moniker.String)
	suite.Require().True(profile.CreationDate.Equal(stored.CreationDate.Time))

	// Override data
	changedProfile := profile.
		WithMoniker(newStrPtr("second-moniker")).
		WithBio(newStrPtr("biography")).
		WithPictures(nil, newStrPtr("cover-picture"))

	err = suite.database.SaveProfile(changedProfile)
	suite.Require().NoError(err, "overriding profile should return no error")

	// Verify the storing
	overridden, err := suite.database.GetUserByAddress(profile.Creator)
	suite.Require().NoError(err)

	suite.Require().Equal(profile.DTag, overridden.DTag.String)
	suite.Require().Equal(*changedProfile.Moniker, overridden.Moniker.String)
	suite.Require().Equal(*changedProfile.Bio, overridden.Bio.String)
	suite.Require().False(overridden.ProfilePic.Valid)
	suite.Require().Equal(*changedProfile.Pictures.Cover, overridden.CoverPic.String)
	suite.Require().True(profile.CreationDate.Equal(overridden.CreationDate.Time))
}
