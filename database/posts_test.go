package database_test

import (
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
)

func (suite *DbTestSuite) TestDesmosDb_SavePost() {
	created, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	post := poststypes.NewPost(
		"979cc7397c87be773dd04fd219cdc031482efc9ed5443b7b636de1aff0179fc4",
		"",
		"Post message",
		false,
		"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
		poststypes.OptionalData{
			poststypes.NewOptionalDataEntry("first_key", "first_value"),
			poststypes.NewOptionalDataEntry("second_key", "1"),
		},
		poststypes.NewAttachments(
			poststypes.NewAttachment(
				"http://example.com/uri",
				"image/png",
				[]string{
					"cosmos1h7snyfa2kqyea2kelnywzlmle9vfmj3378xfkn",
					"cosmos19aa4ys9vy98unh68r6hc2sqhgv6ze4svrxh2vn",
				},
			),
		),
		poststypes.NewPollData(
			"Do you like dogs?",
			time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			[]poststypes.PollAnswer{
				poststypes.NewPollAnswer("1", "Yes"),
				poststypes.NewPollAnswer("2", "No"),
			},
			true,
			false,
		),
		time.Time{},
		created,
		"cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
	)

	// Save the data
	err = suite.database.SavePost(post)
	suite.Require().NoError(err)

	// Get the data
	stored, err := suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)
	suite.Require().NotNil(stored)
	suite.Require().True(post.Equal(stored))
}
