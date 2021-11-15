package database_test

import (
	"time"

	"github.com/desmos-labs/djuno/types"

	dbtypes "github.com/desmos-labs/djuno/database/types"

	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"
)

func (suite *DbTestSuite) TestDesmosDb_SavePost() {
	err := suite.database.SaveUserIfNotExisting("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d", 1)
	suite.Require().NoError(err)

	post := types.NewPost(
		poststypes.NewPost(
			"979cc7397c87be773dd04fd219cdc031482efc9ed5443b7b636de1aff0179fc4",
			"",
			"Post message",
			poststypes.CommentsStateBlocked,
			"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			[]poststypes.Attribute{
				poststypes.NewAttribute("first_key", "first_value"),
				poststypes.NewAttribute("second_key", "1"),
			},
			poststypes.NewAttachments(
				poststypes.NewAttachment(
					"https://example.com/uri",
					"image/png",
					[]string{
						"cosmos1h7snyfa2kqyea2kelnywzlmle9vfmj3378xfkn",
						"cosmos19aa4ys9vy98unh68r6hc2sqhgv6ze4svrxh2vn",
					},
				),
			),
			poststypes.NewPoll(
				"Do you like dogs?",
				time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				[]poststypes.ProvidedAnswer{
					poststypes.NewProvidedAnswer("1", "Yes"),
					poststypes.NewProvidedAnswer("2", "No"),
				},
				true,
				false,
			),
			time.Time{},
			time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
			"cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
		),
		10,
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

func (suite *DbTestSuite) savePollData() (poststypes.Post, *poststypes.Poll) {
	post := suite.testData.post
	err := suite.database.SaveUserIfNotExisting(post.Creator, 1)
	suite.Require().NoError(err)

	err = suite.database.SavePost(types.NewPost(post, 1))
	suite.Require().NoError(err)

	return post, post.Poll
}

func (suite *DbTestSuite) TestDesmosDb_SavePollAnswer() {
	post, _ := suite.savePollData()

	// Save the answer
	user := "cosmos184dqecwkwex2hv6ae8fhzkw0cwrn39aw2ncy7n"
	err := suite.database.SaveUserIfNotExisting(user, 1)
	suite.Require().NoError(err)

	err = suite.database.SaveUserPollAnswer(types.NewUserPollAnswer(
		poststypes.NewUserAnswer(post.PostID, user, []string{"0", "1"}),
		1,
	))
	suite.Require().NoError(err)

	// Verify the insertion
	var rows []dbtypes.UserPollAnswerRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_poll_answer")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 2)
	suite.Require().True(rows[0].Equal(dbtypes.NewUserPollAnswerRow(1, "0", user, 1)))
	suite.Require().True(rows[1].Equal(dbtypes.NewUserPollAnswerRow(1, "1", user, 1)))
}
