package database_test

import (
	poststypes "github.com/desmos-labs/desmos/x/posts"
	"time"
)

func (suite *DbTestSuite) savePost() poststypes.Post {
	post := suite.testData.post

	err := suite.database.SavePost(post)
	suite.Require().NoError(err)

	return post
}

func (suite *DbTestSuite) TestDesmosDb_SavePost() {
	// Save the data
	post := suite.savePost()

	// Get the data
	stored, err := suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)
	suite.NotNil(stored)
	suite.Require().True(post.Equals(*stored))

	// Trying to replace does nothing
	post.Message = "Edited post message"
	err = suite.database.SavePost(post)
	suite.Require().NoError(err)

	// Get the data
	stored, err = suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)
	suite.NotNil(stored)
	suite.NotEqual(post.Message, stored.Message)
}

func (suite *DbTestSuite) TestDesmosDb_EditPost() {
	// Save the post
	post := suite.savePost()

	// Edit the post
	editDate, err := time.Parse(time.RFC3339, "2020-12-31T00:00:00Z")
	suite.Require().NoError(err)

	err = suite.database.EditPost(post.PostID, post.Message+"-edited", editDate)
	suite.Require().NoError(err)

	// Check the edit
	stored, err := suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)
	suite.Require().Equal(post.Message+"-edited", stored.Message)
	suite.Require().True(stored.LastEdited.Equal(editDate))
}

func (suite *DbTestSuite) TestDesmosDb_GetPostByID() {
	post := suite.savePost()

	// Verify saving
	stored, err := suite.database.GetPostByID(post.PostID)
	suite.Require().NoError(err)

	suite.True(stored.Equals(post))
}
