package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts"
)

func (suite *DbTestSuite) TestDesmosDb_SavePost() {
	creator, err := sdk.AccAddressFromBech32("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	created, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	post := poststypes.NewPost(
		"979cc7397c87be773dd04fd219cdc031482efc9ed5443b7b636de1aff0179fc4",
		"",
		"Post message",
		false,
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		map[string]string{
			"first_key":  "first_value",
			"second_key": "1",
		},
		created,
		creator,
	)

	// Save the data
	err = suite.database.SavePost(post)
	suite.Require().NoError(err)

	// Get the data
	stored, err := suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)
	suite.Require().NotNil(stored)
	suite.Require().True(post.Equals(*stored))
}
