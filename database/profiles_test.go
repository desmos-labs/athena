package database_test

import (
	"encoding/hex"
	"time"

	relationshipstypes "github.com/desmos-labs/desmos/v4/x/relationships/types"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	"github.com/desmos-labs/djuno/v2/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	dbtypes "github.com/desmos-labs/djuno/v2/database/types"

	profilestypes "github.com/desmos-labs/desmos/v4/x/profiles/types"
)

func (suite *DbTestSuite) TestDesmosDb_SaveUserIfNotExisting() {
	err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
	suite.Require().NoError(err, "storing of address should return no error")

	err = suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
	suite.Require().NoError(err, "storing address second time should return no error")

	user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.GetAddress().String())
}

func (suite *DbTestSuite) verifyEqual(expected, actual *profilestypes.Profile) {
	suite.Require().Equal(expected.Account, actual.Account)
	suite.Require().Equal(expected.DTag, actual.DTag)
	suite.Require().Equal(expected.Nickname, actual.Nickname)
	suite.Require().Equal(expected.Bio, actual.Bio)
	suite.Require().Equal(expected.Pictures, actual.Pictures)
	suite.Require().True(expected.CreationDate.Equal(actual.CreationDate))
}

func (suite *DbTestSuite) TestDesmosDb_SaveProfile() {
	addr, err := sdk.AccAddressFromBech32("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)

	original, err := profilestypes.NewProfile(
		"original-moniker",
		"",
		"",
		profilestypes.NewPictures("", ""),
		time.Time{},
		authtypes.NewBaseAccountWithAddress(addr),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(original, 10))
	suite.Require().NoError(err)

	// Verify the storing
	stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(original, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with a lower height
	updated, err := original.Update(profilestypes.NewProfileUpdate(
		"new-dtag",
		"new-moniker",
		"new-bio",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 9))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(original, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with same height
	updated, err = original.Update(profilestypes.NewProfileUpdate(
		"new-dtag",
		"new-moniker",
		"new-bio",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 10))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(updated, stored)

	// ----------------------------------------------------------------------------------------------------------------

	// Try updating with higher height
	updated, err = original.Update(profilestypes.NewProfileUpdate(
		"new-dtag-2",
		"new-moniker-2",
		"new-bio-2",
		profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify)),
	)
	suite.Require().NoError(err)

	// Save the data
	err = suite.database.SaveProfile(types.NewProfile(updated, 11))
	suite.Require().NoError(err)

	// Verify the data
	stored, err = suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
	suite.Require().NoError(err)
	suite.verifyEqual(updated, stored)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveRelationship() types.Relationship {
	err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
	suite.Require().NoError(err)

	err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
	suite.Require().NoError(err)

	relationship := types.NewRelationship(
		relationshipstypes.NewRelationship(
			"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
			"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
			0,
		),
		10,
	)

	// Save the relationship
	err = suite.database.SaveRelationship(relationship)
	suite.Require().NoError(err)

	return relationship
}

func (suite *DbTestSuite) TestDesmosDb_SaveRelationship() {
	relationship := suite.saveRelationship()

	err := suite.database.SaveRelationship(relationship)
	suite.Require().NoError(err, "double inserting the same relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM profile_relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewRelationshipRow(
		relationship.Creator,
		relationship.Counterparty,
		relationship.SubspaceID,
		relationship.Height,
	)))
}

func (suite *DbTestSuite) TestDesmosDb_DeleteRelationship() {
	relationship := suite.saveRelationship()

	err := suite.database.DeleteRelationship(relationship)
	suite.Require().NoError(err, "removing existing relationship should return no error")

	var rows []dbtypes.RelationshipRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM profile_relationship")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.DeleteRelationship(relationship)
	suite.Require().NoError(err, "deleting non existent relationship should return no error")
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveBlockage() types.Blockage {
	suite.Require().NoError(suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1))
	suite.Require().NoError(suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1))

	blockage := types.NewBlockage(
		relationshipstypes.NewUserBlock(
			"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
			"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
			"this is my blocking reason",
			0,
		),
		1,
	)

	// Save the blockage
	err := suite.database.SaveUserBlock(blockage)
	suite.Require().NoError(err)

	return blockage
}

func (suite *DbTestSuite) TestDesmosDB_SaveUserBlockage() {
	blockage := suite.saveBlockage()

	err := suite.database.SaveUserBlock(blockage)
	suite.Require().NoError(err, "double inserting blockage should return no error")

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.NewBlockageRow(
		blockage.Blocker,
		blockage.Blocked,
		blockage.Reason,
		blockage.SubspaceID,
		blockage.Height,
	)))
}

func (suite *DbTestSuite) TestDesmosDB_RemoveUserBlockage() {
	blockage := suite.saveBlockage()

	err := suite.database.DeleteBlockage(blockage)
	suite.Require().NoError(err)

	var rows []dbtypes.BlockageRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_block")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 0)

	err = suite.database.DeleteBlockage(blockage)
	suite.Require().NoError(err, "deleting non existing blockage should return no error")
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestDesmosDB_SaveChainLink() {
	bz, err := sdk.GetFromBech32("desmospub1addwnpepqvczf60q448expz77knqhwfpw8nyrrx38vyzu7nmc0ks2vf2pdqh63tcdmy", "desmospub")
	suite.Require().NoError(err)

	pubKey, err := legacy.PubKeyFromBytes(bz)
	suite.Require().NoError(err)

	signature, err := hex.DecodeString("74657874")
	suite.Require().NoError(err)

	chainLink := types.NewChainLink(
		profilestypes.NewChainLink(
			"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
			profilestypes.NewBech32Address("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "desmos"),
			profilestypes.NewProof(pubKey, &profilestypes.SingleSignatureData{Mode: 1, Signature: signature}, "text"),
			profilestypes.NewChainConfig("desmos"),
			time.Now(),
		),
		10,
	)
	err = suite.database.SaveUserIfNotExisting("cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f", 10)
	suite.Require().NoError(err)

	err = suite.database.SaveChainLink(chainLink)
	suite.Require().NoError(err)
}

func (suite *DbTestSuite) TestDesmosDB_DeleteChainLink() {
	bz, err := sdk.GetFromBech32("desmospub1addwnpepqvczf60q448expz77knqhwfpw8nyrrx38vyzu7nmc0ks2vf2pdqh63tcdmy", "desmospub")
	suite.Require().NoError(err)

	pubKey, err := legacy.PubKeyFromBytes(bz)
	suite.Require().NoError(err)

	signature, err := hex.DecodeString("74657874")
	suite.Require().NoError(err)

	chainLink := types.NewChainLink(
		profilestypes.NewChainLink(
			"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
			profilestypes.NewBech32Address("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "desmos"),
			profilestypes.NewProof(pubKey, &profilestypes.SingleSignatureData{Mode: 1, Signature: signature}, "text"),
			profilestypes.NewChainConfig("desmos"),
			time.Now(),
		),
		10,
	)
	err = suite.database.SaveUserIfNotExisting("cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f", 10)
	suite.Require().NoError(err)

	err = suite.database.SaveChainLink(chainLink)
	suite.Require().NoError(err)

	err = suite.database.DeleteChainLink(
		"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
		"desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
		"desmos",
		10,
	)
	suite.Require().NoError(err)

	var count int
	err = suite.database.Sql.QueryRow("SELECT COUNT(id) FROM chain_link").Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Zero(count)

	err = suite.database.Sql.QueryRow("SELECT COUNT(id) FROM chain_link_proof").Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Zero(count)

	err = suite.database.Sql.QueryRow("SELECT COUNT(id) FROM chain_link_chain_config").Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestDesmosDB_DeleteApplicationLink() {
	user := "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47"
	err := suite.database.SaveUserIfNotExisting(user, 1)
	suite.Require().NoError(err)

	applicationLink := types.NewApplicationLink(
		profilestypes.NewApplicationLink(
			user,
			profilestypes.NewData("twitter", "twitteruser"),
			profilestypes.ApplicationLinkStateInitialized,
			profilestypes.NewOracleRequest(
				0,
				1,
				profilestypes.NewOracleRequestCallData(
					"twitter",
					"7B22757365726E616D65223A22526963636172646F4D222C22676973745F6964223A223732306530303732333930613930316262383065353966643630643766646564227D",
				),
				"client_id",
			),
			nil,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		),
		100,
	)

	err = suite.database.SaveApplicationLink(applicationLink)
	suite.Require().NoError(err)

	var count int
	err = suite.database.Sql.QueryRow("SELECT COUNT(*) FROM application_link").Scan(&count)
	suite.Require().NoError(err)
	suite.Require().Equal(1, count)

	err = suite.database.DeleteApplicationLink(user, "twitter", "twitteruser", 100)
	suite.Require().NoError(err)
}
