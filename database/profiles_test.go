package database_test

import (
	"time"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err, "storing of address should return no error")

	err = suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err, "storing address second time should return no error")

	user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.Creator)
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	profile := profilestypes.NewProfile(
		"profile-moniker",
		"",
		"",
		profilestypes.NewPictures("", ""),
		time.Time{},
		"cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f",
	)

	// Save data
	err := suite.database.SaveProfile(profile)
	suite.Require().NoError(err)

	// Override data
	changedProfile := profilestypes.NewProfile(
		"second-moniker",
		"",
		"biography",
		profilestypes.NewPictures("", "cover-picture"),
		time.Time{},
		"cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f",
	)

	err = suite.database.SaveProfile(changedProfile)
	suite.Require().NoError(err, "overriding profile should return no error")

	// Verify the storing
	stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.Require().True(changedProfile.Equal(stored))
}
