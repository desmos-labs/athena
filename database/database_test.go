package database_test

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	poststypes "github.com/desmos-labs/desmos/x/staging/posts/types"

	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/djuno/database"
	jconfig "github.com/desmos-labs/juno/config"
	"github.com/stretchr/testify/suite"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database *database.DesmosDb
	testData TestData
}

type TestData struct {
	post poststypes.Post
}

func (suite *DbTestSuite) SetupTest() {
	// Setup test data
	suite.setupTestData()

	// Create the codec
	_, codec := desmosapp.MakeCodecs()

	// Build the database
	config := &jconfig.Config{
		DatabaseConfig: &jconfig.DatabaseConfig{
			Type: "psql",
			Config: &jconfig.PostgreSQLConfig{
				Name:     "juno",
				Host:     "localhost",
				Port:     5433,
				User:     "juno",
				Password: "password",
			},
		},
	}

	db, err := database.Builder(config, codec)
	suite.Require().NoError(err)

	desmosDb, ok := (db).(*database.DesmosDb)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = desmosDb.Sql.Exec(`DROP SCHEMA public CASCADE;`)
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = desmosDb.Sql.Exec(`CREATE SCHEMA public;`)
	suite.Require().NoError(err)

	dirPath := "../schema"
	dir, err := ioutil.ReadDir(dirPath)
	for _, fileInfo := range dir {
		if !strings.HasSuffix(fileInfo.Name(), ".sql") {
			continue
		}

		file, err := ioutil.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := desmosDb.Sql.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	suite.database = desmosDb
}

func (suite *DbTestSuite) setupTestData() {
	suite.testData = TestData{
		post: poststypes.NewPost(
			"60303ae22b998861bce3b28f33eec1be758a213c86c93c076dbe9f558c11c752",
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
			time.Date(2020, 10, 10, 15, 00, 00, 00, time.UTC),
			"cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
		),
	}
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
