package database_test

import (
	"time"

	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	dbtypes "github.com/desmos-labs/djuno/database/types"
)

func (suite *DbTestSuite) savePollData() (poststypes.Post, *poststypes.PollData) {
	// Save the post
	post := suite.testData.post
	err := suite.database.SavePost(post)
	suite.Require().NoError(err)

	// Save the pollData
	pollData := poststypes.NewPollData(
		"Which are better: dogs or cats?",
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		poststypes.NewPollAnswers(
			poststypes.NewPollAnswer("0", "Cats"),
			poststypes.NewPollAnswer("1", "Dogs"),
		),
		false,
		false,
	)

	err = suite.database.SavePollData(post.PostID, pollData)
	suite.Require().NoError(err)

	return post, pollData
}

func (suite *DbTestSuite) TestDesmosDb_SavePollData() {
	post, pollData := suite.savePollData()

	// Get the inserted data
	var rows []dbtypes.PollRow
	stmt := `SELECT * FROM poll WHERE post_id = $1`
	err := suite.database.Sqlx.Select(&rows, stmt, post.PostID)
	suite.Require().NoError(err)

	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equal(dbtypes.PollRow{
		Id:                    1,
		PostID:                post.PostID,
		Question:              pollData.Question,
		EndDate:               pollData.EndDate,
		AllowsMultipleAnswers: pollData.AllowsMultipleAnswers,
		AllowsAnswerEdits:     pollData.AllowsAnswerEdits,
	}))
}

func (suite *DbTestSuite) TestDesmosDb_SavePollAnswer() {
	post, _ := suite.savePollData()

	// Save the answer
	user := "cosmos184dqecwkwex2hv6ae8fhzkw0cwrn39aw2ncy7n"
	answer := poststypes.NewUserAnswer([]string{"0", "1"}, user)
	err := suite.database.SaveUserPollAnswer(post.PostID, answer)
	suite.Require().NoError(err)

	// Verify the insertion
	var rows []dbtypes.UserPollAnswerRow
	err = suite.database.Sqlx.Select(&rows, "SELECT * FROM user_poll_answer")
	suite.Require().NoError(err)

	suite.Require().Len(rows, 2)
	suite.Require().True(rows[0].Equal(dbtypes.UserPollAnswerRow{
		PollID:          1,
		Answer:          "0",
		AnswererAddress: user,
	}))
	suite.Require().True(rows[1].Equal(dbtypes.UserPollAnswerRow{
		PollID:          1,
		Answer:          "1",
		AnswererAddress: user,
	}))
}

func (suite *DbTestSuite) TestDesmosDb_GetPollByPostId() {
	post, poll := suite.savePollData()

	// Get the data
	stored, err := suite.database.GetPollByPostID(post.PostID)
	suite.Require().NoError(err)

	expected := dbtypes.PollRow{
		Id:                    1,
		PostID:                post.PostID,
		Question:              poll.Question,
		EndDate:               poll.EndDate,
		AllowsMultipleAnswers: poll.AllowsMultipleAnswers,
		AllowsAnswerEdits:     poll.AllowsAnswerEdits,
	}
	suite.Require().True(stored.Equal(expected))
}
