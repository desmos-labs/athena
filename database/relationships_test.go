package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"
)

func (suite *DbTestSuite) saveRelationship() (sender, receiver sdk.AccAddress, subspace string) {
	sender, err := sdk.AccAddressFromBech32("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
	suite.Require().NoError(err)

	receiver, err = sdk.AccAddressFromBech32("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr")
	suite.Require().NoError(err)

	subspace = "mooncake"

	// Save the relationship
	err = suite.database.SaveRelationship(sender, receiver, subspace)
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
		Sender:   sender.String(),
		Receiver: receiver.String(),
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

// ____________________________________

func (suite *DbTestSuite) saveBlockage() (blocker, blocked sdk.AccAddress, reason, subspace string) {
	blocker, err := sdk.AccAddressFromBech32("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
	suite.Require().NoError(err)

	blocked, err = sdk.AccAddressFromBech32("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr")
	suite.Require().NoError(err)

	reason = "this is my blocking reason"
	subspace = "mooncake"

	// Save the blockage
	err = suite.database.SaveBlockage(blocker, blocked, reason, subspace)
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
		Blocker:  blocker.String(),
		Blocked:  blocked.String(),
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
