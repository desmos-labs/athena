package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles"
)

func newStrPtr(value string) *string {
	return &value
}

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	addr, err := sdk.AccAddressFromBech32("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	user, err := suite.database.SaveUserIfNotExisting(addr)
	suite.Require().NoError(err, "storing of address should return no error")

	otherUser, err := suite.database.SaveUserIfNotExisting(addr)
	suite.Require().NoError(err, "storing address second time should return no error")
	suite.Require().Equal(user, otherUser)

	storedAddr, err := sdk.AccAddressFromBech32(user.Address)
	suite.Require().NoError(err)
	suite.Require().True(storedAddr.Equals(addr))
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	creator, err := sdk.AccAddressFromBech32("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)

	profile := profilestypes.NewProfile(creator).
		WithMoniker("profile-moniker")

	// Save data
	stored, err := suite.database.SaveProfile(profile)
	suite.Require().NoError(err)

	moniker := stored.Moniker.String
	suite.Require().Equal(profile.Moniker, moniker)

	// Override data
	changedProfile := profile.
		WithMoniker("second-moniker").
		WithBio(newStrPtr("biography")).
		WithName(newStrPtr("custom-name")).
		WithSurname(newStrPtr("custom-surname")).
		WithPictures(nil, newStrPtr("cover-picture"))

	overridden, err := suite.database.SaveProfile(changedProfile)
	suite.Require().NoError(err, "overriding profile should return no error")

	// Verify the storing
	newMoniker := overridden.Moniker.String
	suite.Require().Equal(changedProfile.Moniker, newMoniker)

	newBio := overridden.Bio.String
	suite.Require().Equal(*changedProfile.Bio, newBio)

	newName := overridden.Name.String
	suite.Require().Equal(*changedProfile.Name, newName)

	newSurname := overridden.Surname.String
	suite.Require().Equal(*changedProfile.Surname, newSurname)

	suite.Require().False(overridden.ProfilePic.Valid)

	newCover := overridden.CoverPic.String
	suite.Require().Equal(*changedProfile.Pictures.Cover, newCover)
}
