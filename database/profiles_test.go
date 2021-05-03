package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
)

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err, "storing of address should return no error")

	err = suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err, "storing address second time should return no error")

	user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.GetAddress().String())
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	addr1, err := sdk.AccAddressFromBech32("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)

	profile, err := profilestypes.NewProfile(
		"profile-moniker",
		"",
		"",
		profilestypes.NewPictures("", ""),
		time.Time{},
		authtypes.NewBaseAccountWithAddress(addr1),
	)
	suite.Require().NoError(err)

	// Save data
	err = suite.database.SaveProfile(profile)
	suite.Require().NoError(err)

	// Override data
	changedProfile, err := profilestypes.NewProfile(
		"second-moniker",
		"",
		"biography",
		profilestypes.NewPictures("", "cover-picture"),
		time.Time{},
		authtypes.NewBaseAccountWithAddress(addr1),
	)

	err = suite.database.SaveProfile(changedProfile)
	suite.Require().NoError(err, "overriding profile should return no error")

	// Verify the storing
	stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.Require().Equal(changedProfile, stored)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveRelationship() (sender, receiver string, subspace string) {
	sender = "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s"
	receiver = "cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr"
	subspace = "mooncake"

	// Save the relationship
	err := suite.database.SaveRelationship(sender, receiver, subspace)
	suite.Require().NoError(err)

	return sender, receiver, subspace
}

func (suite *DbTestSuite) TestDesmosDb_SaveRelationship() {
	sender, receiver, subspace := suite.saveRelationship()

	err := suite.database.SaveRelationship(sender, receiver, subspace)
	suite.Require().NoError(err, "double inserting the same relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.RelationshipRow{
		Sender:   sender,
		Receiver: receiver,
		Subspace: subspace,
	}))
}

func (suite *DbTestSuite) TestDesmosDb_DeleteRelationship() {
	sender, receiver, subspace := suite.saveRelationship()

	err := suite.database.DeleteRelationship(sender, receiver, subspace)
	suite.Require().NoError(err, "removing existing relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.DeleteRelationship(sender, receiver, subspace)
	suite.Require().NoError(err, "deleting non existent relationship should return no error")
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveBlockage() (blocker, blocked string, reason, subspace string) {
	blocker = "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s"
	blocked = "cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr"
	reason = "this is my blocking reason"
	subspace = "mooncake"

	// Save the blockage
	err := suite.database.SaveBlockage(blocker, blocked, reason, subspace)
	suite.Require().NoError(err)

	return blocker, blocked, reason, subspace
}

func (suite *DbTestSuite) TestDesmosDB_SaveUserBlockage() {
	blocker, blocked, reason, subspace := suite.saveBlockage()

	err := suite.database.SaveBlockage(blocker, blocked, reason, subspace)
	suite.Require().NoError(err, "double inserting blockage should return no error")

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.BlockageRow{
		Blocker:  blocker,
		Blocked:  blocked,
		Reason:   reason,
		Subspace: subspace,
	}))
}

func (suite *DbTestSuite) TestDesmosDB_RemoveUserBlockage() {
	blocker, blocked, _, subspace := suite.saveBlockage()

	err := suite.database.RemoveBlockage(blocker, blocked, subspace)
	suite.Require().NoError(err)

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.RemoveBlockage(blocker, blocked, subspace)
	suite.Require().NoError(err, "deleting non existing blockage should return no error")
}
