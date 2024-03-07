package database_test

import (
	"encoding/hex"
	"time"

	"github.com/cosmos/cosmos-sdk/codec/legacy"

	"github.com/desmos-labs/athena/v2/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/v7/x/profiles/types"
)

func (suite *DbTestSuite) TestSaveUserIfNotExisting() {
	testCases := []struct {
		name      string
		setup     func()
		address   string
		shouldErr bool
		check     func()
	}{
		{
			name:      "non existing user returns no error",
			address:   "cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
			shouldErr: false,
			check: func() {
				user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
				suite.Require().NoError(err)
				suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.GetAddress().String())
			},
		},
		{
			name: "existing user returns no error",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
				suite.Require().NoError(err)
			},
			address:   "cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
			shouldErr: false,
			check: func() {
				user, err := suite.database.GetUserByAddress("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
				suite.Require().NoError(err)
				suite.Require().Equal("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", user.GetAddress().String())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveUserIfNotExisting(tc.address, 1)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) verifyEqual(expected, actual *profilestypes.Profile) {
	suite.Require().Equal(expected.Account, actual.Account)
	suite.Require().Equal(expected.DTag, actual.DTag)
	suite.Require().Equal(expected.Nickname, actual.Nickname)
	suite.Require().Equal(expected.Bio, actual.Bio)
	suite.Require().Equal(expected.Pictures, actual.Pictures)
	suite.Require().True(expected.CreationDate.Equal(actual.CreationDate))
}

func (suite *DbTestSuite) updateProfile(profile *profilestypes.Profile, update *profilestypes.ProfileUpdate) *profilestypes.Profile {
	updated, err := profile.Update(update)
	suite.Require().NoError(err)
	return updated
}

func (suite *DbTestSuite) TestSaveProfile() {
	testCases := []struct {
		name      string
		setup     func()
		profile   *profilestypes.Profile
		height    int64
		shouldErr bool
		check     func()
	}{
		{
			name:      "non existing profile is stored properly",
			profile:   suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f"),
			height:    10,
			shouldErr: false,
			check: func() {
				stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				suite.Require().NoError(err)
				suite.verifyEqual(suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f"), stored)
			},
		},
		{
			name: "updating with a lower height does nothing",
			setup: func() {
				profile := suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				err := suite.database.SaveProfile(types.NewProfile(profile, 10))
				suite.Require().NoError(err)
			},
			profile: suite.updateProfile(
				suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f"),
				profilestypes.NewProfileUpdate(
					"NewDTag",
					"New Nickname",
					"This is my new biography",
					profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify),
				),
			),
			height:    9,
			shouldErr: false,
			check: func() {
				stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				suite.Require().NoError(err)

				// Make sure the profile has not been updated
				suite.Require().Equal("TestUser", stored.DTag)
				suite.Require().Equal("Test User", stored.Nickname)
				suite.Require().Equal("This is a test user", stored.Bio)
			},
		},
		{
			name: "updating with the same height updates the profile",
			setup: func() {
				profile := suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				err := suite.database.SaveProfile(types.NewProfile(profile, 10))
				suite.Require().NoError(err)
			},
			profile: suite.updateProfile(
				suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f"),
				profilestypes.NewProfileUpdate(
					"NewDTag",
					"New Nickname",
					"This is my new biography",
					profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify),
				),
			),
			height: 10,
			check: func() {
				stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				suite.Require().NoError(err)

				// Make sure the profile has been updated
				suite.Require().Equal("NewDTag", stored.DTag)
				suite.Require().Equal("New Nickname", stored.Nickname)
				suite.Require().Equal("This is my new biography", stored.Bio)
			},
		},
		{
			name: "updating with a higher height updates the profile",
			setup: func() {
				profile := suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				err := suite.database.SaveProfile(types.NewProfile(profile, 10))
				suite.Require().NoError(err)
			},
			profile: suite.updateProfile(
				suite.buildProfile("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f"),
				profilestypes.NewProfileUpdate(
					"NewDTag",
					"New Nickname",
					"This is my new biography",
					profilestypes.NewPictures(profilestypes.DoNotModify, profilestypes.DoNotModify),
				),
			),
			height: 11,
			check: func() {
				stored, err := suite.database.GetUserByAddress("cosmos15c66kjz44zm58xqlcqjwftan4tnaeq7rtmhn4f")
				suite.Require().NoError(err)

				// Make sure the profile has been updated
				suite.Require().Equal("NewDTag", stored.DTag)
				suite.Require().Equal("New Nickname", stored.Nickname)
				suite.Require().Equal("This is my new biography", stored.Bio)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveProfile(types.NewProfile(tc.profile, tc.height))
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestSaveChainLink() {
	testCases := []struct {
		name           string
		setup          func()
		buildChainLink func() types.ChainLink
		shouldErr      bool
		check          func()
	}{
		{
			name: "non existing chain link is stored properly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f", 10)
				suite.Require().NoError(err)
			},
			buildChainLink: func() types.ChainLink {
				bz, err := sdk.GetFromBech32("desmospub1addwnpepqvczf60q448expz77knqhwfpw8nyrrx38vyzu7nmc0ks2vf2pdqh63tcdmy", "desmospub")
				suite.Require().NoError(err)

				pubKey, err := legacy.PubKeyFromBytes(bz)
				suite.Require().NoError(err)

				signature, err := hex.DecodeString("74657874")
				suite.Require().NoError(err)

				return types.NewChainLink(
					profilestypes.NewChainLink(
						"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
						profilestypes.NewBech32Address("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "desmos"),
						profilestypes.NewProof(pubKey, &profilestypes.SingleSignature{ValueType: 1, Signature: signature}, "text"),
						profilestypes.NewChainConfig("desmos"),
						time.Now(),
					),
					10,
				)
			},
			shouldErr: false,
			check: func() {
				// Make sure the chain link has been stored
				var count int
				err := suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)

				err = suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link_proof").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)

				err = suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link_chain_config").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)

				// Make sure the chain link count has been updated
				var chainLinkCount int
				err = suite.database.SQL.Get(&chainLinkCount, "SELECT chain_links_count FROM profile_counters WHERE profile_address = $1", "cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f")
				suite.Require().NoError(err)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveChainLink(tc.buildChainLink())
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}

func (suite *DbTestSuite) TestDeleteChainLink() {
	testCases := []struct {
		name            string
		setup           func()
		user            string
		externalAddress string
		chainName       string
		height          int64
		shouldErr       bool
		check           func()
	}{
		{
			name: "existing chain link is deleted properly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f", 10)
				suite.Require().NoError(err)

				bz, err := sdk.GetFromBech32("desmospub1addwnpepqvczf60q448expz77knqhwfpw8nyrrx38vyzu7nmc0ks2vf2pdqh63tcdmy", "desmospub")
				suite.Require().NoError(err)

				pubKey, err := legacy.PubKeyFromBytes(bz)
				suite.Require().NoError(err)

				signature, err := hex.DecodeString("74657874")
				suite.Require().NoError(err)

				err = suite.database.SaveChainLink(types.NewChainLink(
					profilestypes.NewChainLink(
						"cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
						profilestypes.NewBech32Address("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "desmos"),
						profilestypes.NewProof(pubKey, &profilestypes.SingleSignature{ValueType: 1, Signature: signature}, "text"),
						profilestypes.NewChainConfig("desmos"),
						time.Now(),
					),
					10,
				))
				suite.Require().NoError(err)
			},
			user:            "cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f",
			externalAddress: "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
			chainName:       "desmos",
			height:          10,
			shouldErr:       false,
			check: func() {
				// Make sure the chain link has been deleted
				var count int
				err := suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Zero(count)

				err = suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link_proof").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Zero(count)

				// The chain config should not be deleted as this is something that is shared among all the chain links
				err = suite.database.SQL.QueryRow("SELECT COUNT(id) FROM chain_link_chain_config").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)

				// Make sure the chain link count has been updated
				var chainLinkCount int
				err = suite.database.SQL.Get(&chainLinkCount, "SELECT chain_links_count FROM profile_counters WHERE profile_address = $1", "cosmos10clxpupsmddtj7wu7g0wdysajqwp890mva046f")
				suite.Require().NoError(err)
				suite.Require().Zero(chainLinkCount)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.DeleteChainLink(tc.user, tc.externalAddress, tc.chainName, tc.height)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestSaveApplicationLink() {
	testCase := []struct {
		name            string
		setup           func()
		applicationLink types.ApplicationLink
		shouldErr       bool
		check           func()
	}{
		{
			name: "non existing application link is stored properly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47", 1)
				suite.Require().NoError(err)
			},
			applicationLink: types.NewApplicationLink(
				profilestypes.NewApplicationLink(
					"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
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
					time.Date(2021, 1, 1, 00, 00, 00, 000, time.UTC),
				),
				100,
			),
			shouldErr: false,
			check: func() {
				// Make sure the application link has been stored
				var count int
				err := suite.database.SQL.QueryRow("SELECT COUNT(id) FROM application_link").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)

				// Make sure the application link count has been updated
				var applicationLinkCount int
				err = suite.database.SQL.Get(&applicationLinkCount, "SELECT application_links_count FROM profile_counters WHERE profile_address = $1", "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
				suite.Require().NoError(err)
				suite.Require().Equal(1, applicationLinkCount)
			},
		},
	}

	for _, tc := range testCase {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveApplicationLink(tc.applicationLink)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}

func (suite *DbTestSuite) TestDeleteApplicationLink() {
	testCases := []struct {
		name        string
		setup       func()
		user        string
		application string
		username    string
		height      int64
		shouldErr   bool
		check       func()
	}{
		{
			name: "existing application link is deleted properly",
			setup: func() {
				err := suite.database.SaveUserIfNotExisting("cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47", 1)
				suite.Require().NoError(err)

				err = suite.database.SaveApplicationLink(types.NewApplicationLink(
					profilestypes.NewApplicationLink(
						"cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
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
						time.Date(2021, 1, 1, 00, 00, 00, 000, time.UTC),
					),
					100,
				))
				suite.Require().NoError(err)
			},
			user:        "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47",
			application: "twitter",
			username:    "twitteruser",
			height:      100,
			shouldErr:   false,
			check: func() {
				// Make sure the application link has been deleted
				var count int
				err := suite.database.SQL.QueryRow("SELECT COUNT(id) FROM application_link").Scan(&count)
				suite.Require().NoError(err)
				suite.Require().Zero(count)

				// Make sure the application link count has been updated
				var applicationLinkCount int
				err = suite.database.SQL.Get(&applicationLinkCount, "SELECT application_links_count FROM profile_counters WHERE profile_address = $1", "cosmos1y54exmx84cqtasvjnskf9f63djuuj68p7hqf47")
				suite.Require().NoError(err)
				suite.Require().Zero(applicationLinkCount)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.DeleteApplicationLink(tc.user, tc.application, tc.username, tc.height)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				tc.check()
			}
		})
	}
}
