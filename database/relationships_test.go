package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	relationshipstypes "github.com/desmos-labs/desmos/v6/x/relationships/types"
	subspacestypes "github.com/desmos-labs/desmos/v6/x/subspaces/types"

	dbtypes "github.com/desmos-labs/athena/database/types"
	"github.com/desmos-labs/athena/types"
)

func (suite *DbTestSuite) TestSaveRelationship() {
	testCases := []struct {
		name         string
		setup        func()
		relationship types.Relationship
		shouldErr    bool
		check        func()
	}{
		{
			name: "relationship saved correctly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveSubspace(types.NewSubspace(subspacestypes.NewSubspace(
					0,
					"Test subspace",
					"",
					"",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					time.Now(),
					sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100000))),
				), 1))
				suite.Require().NoError(err)
			},
			relationship: types.NewRelationship(
				relationshipstypes.NewRelationship(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					0,
				),
				10,
			),
			shouldErr: false,
			check: func() {
				// Make sure the relationship has been saved
				var rows []dbtypes.RelationshipRow
				err := suite.database.SQL.Select(&rows, "SELECT * FROM user_relationship")
				suite.Require().NoError(err)

				suite.Require().Len(rows, 1)
				suite.Require().True(rows[0].Equal(dbtypes.NewRelationshipRow(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					0,
					10,
				)))

				// Make sure the creator's relationships count has been updated
				var count int
				err = suite.database.SQL.Get(&count, "SELECT relationships_count FROM profile_counters WHERE profile_address = $1", "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveRelationship(tc.relationship)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}

func (suite *DbTestSuite) TestDeleteRelationship() {
	testCases := []struct {
		name         string
		setup        func()
		relationship types.Relationship
		shouldErr    bool
		check        func()
	}{
		{
			name: "relationship deleted correctly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveSubspace(types.NewSubspace(subspacestypes.NewSubspace(
					0,
					"Test subspace",
					"",
					"",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					time.Now(),
					sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(100000))),
				), 1))
				suite.Require().NoError(err)

				err = suite.database.SaveRelationship(types.NewRelationship(
					relationshipstypes.NewRelationship(
						"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
						"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
						0,
					),
					10,
				))
				suite.Require().NoError(err)
			},
			relationship: types.NewRelationship(
				relationshipstypes.NewRelationship(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					0,
				),
				10,
			),
			shouldErr: false,
			check: func() {
				// Make sure the relationship has been deleted
				var count int
				err := suite.database.SQL.Get(&count, "SELECT COUNT(*) FROM user_relationship")
				suite.Require().NoError(err)
				suite.Require().Zero(count)

				// Make sure the creator's relationships count has been updated
				err = suite.database.SQL.Get(&count, "SELECT relationships_count FROM profile_counters WHERE profile_address = $1", "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
				suite.Require().NoError(err)
				suite.Require().Equal(0, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.DeleteRelationship(tc.relationship)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) saveBlockage() types.Blockage {
	err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
	suite.Require().NoError(err)

	err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
	suite.Require().NoError(err)

	err = suite.database.SaveSubspace(types.NewSubspace(subspacestypes.NewSubspace(
		0,
		"Test subspace",
		"",
		"",
		"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
		"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
		time.Now(),
		nil,
	), 1))
	suite.Require().NoError(err)

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
	err = suite.database.SaveUserBlock(blockage)
	suite.Require().NoError(err)

	return blockage
}

func (suite *DbTestSuite) TestSaveUserBlockage() {
	testCase := []struct {
		name      string
		setup     func()
		blockage  types.Blockage
		shouldErr bool
		check     func()
	}{
		{
			name: "blockage saved correctly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveSubspace(types.NewSubspace(subspacestypes.NewSubspace(
					0,
					"Test subspace",
					"",
					"",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					time.Now(),
					nil,
				), 1))
				suite.Require().NoError(err)
			},
			blockage: types.NewBlockage(
				relationshipstypes.NewUserBlock(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					"this is my blocking reason",
					0,
				),
				1,
			),
			shouldErr: false,
			check: func() {
				// Make sure the blockage has been saved
				var rows []dbtypes.BlockageRow
				err := suite.database.SQL.Select(&rows, "SELECT * FROM user_block")
				suite.Require().NoError(err)

				suite.Require().Len(rows, 1)
				suite.Require().True(rows[0].Equal(dbtypes.NewBlockageRow(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					"this is my blocking reason",
					0,
					1,
				)))

				// Make sure the blocker's blocks count has been updated
				var count int
				err = suite.database.SQL.Get(&count, "SELECT blocks_count FROM profile_counters WHERE profile_address = $1", "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCase {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveUserBlock(tc.blockage)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}

func (suite *DbTestSuite) TestDeleteUserBlockage() {
	testCases := []struct {
		name      string
		setup     func()
		blockage  types.Blockage
		shouldErr bool
		check     func()
	}{
		{
			name: "blockage deleted correctly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveUserIfNotExisting("cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveSubspace(types.NewSubspace(subspacestypes.NewSubspace(
					0,
					"Test subspace",
					"",
					"",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					time.Now(),
					nil,
				), 1))
				suite.Require().NoError(err)

				err = suite.database.SaveUserBlock(types.NewBlockage(
					relationshipstypes.NewUserBlock(
						"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
						"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
						"this is my blocking reason",
						0,
					),
					1,
				))
				suite.Require().NoError(err)
			},
			blockage: types.NewBlockage(
				relationshipstypes.NewUserBlock(
					"cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s",
					"cosmos1u0gz4g865yjadxm2hsst388c462agdz7araedr",
					"this is my blocking reason",
					0,
				),
				1,
			),
			shouldErr: false,
			check: func() {
				// Make sure the blockage has been deleted
				var count int
				err := suite.database.SQL.Get(&count, "SELECT COUNT(*) FROM user_block")
				suite.Require().NoError(err)
				suite.Require().Zero(count)

				// Make sure the blocker's blocks count has been updated
				err = suite.database.SQL.Get(&count, "SELECT blocks_count FROM profile_counters WHERE profile_address = $1", "cosmos1jsdja3rsp4lyfup3pc2r05uzusc2e6x3zl285s")
				suite.Require().NoError(err)
				suite.Require().Equal(0, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.DeleteBlockage(tc.blockage)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}
