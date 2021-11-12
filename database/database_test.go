package database_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	junodb "github.com/forbole/juno/v2/database"
	junodbcfg "github.com/forbole/juno/v2/database/config"
	"github.com/forbole/juno/v2/logging"

	poststypes "github.com/desmos-labs/desmos/v2/x/staging/posts/types"

	desmosapp "github.com/desmos-labs/desmos/v2/app"
	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/djuno/database"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database *database.Db
	testData TestData
}

type TestData struct {
	post poststypes.Post
}

func (suite *DbTestSuite) SetupTest() {
	// Setup test data
	suite.setupTestData()

	// Build the database
	encodingConfig := desmosapp.MakeTestEncodingConfig()
	databaseConfig := junodbcfg.NewDatabaseConfig(
		"djuno",
		"localhost",
		5433,
		"djuno",
		"password",
		"",
		"public",
		10,
		10,
	)

	db, err := database.Builder(junodb.NewContext(databaseConfig, &encodingConfig, logging.DefaultLogger()))
	suite.Require().NoError(err)

	desmosDb, ok := (db).(*database.Db)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = desmosDb.Sql.Exec(fmt.Sprintf(`DROP SCHEMA %s CASCADE;`, databaseConfig.Schema))
	suite.Require().NoError(err)

	// Re-create the schema
	_, err = desmosDb.Sql.Exec(fmt.Sprintf(`CREATE SCHEMA %s;`, databaseConfig.Schema))
	suite.Require().NoError(err)

	dirPath := "schema"
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
			poststypes.CommentsStateBlocked,
			"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			[]poststypes.Attribute{
				poststypes.NewAttribute("first_key", "first_value"),
				poststypes.NewAttribute("second_key", "1"),
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
			time.Date(2020, 10, 10, 15, 00, 00, 00, time.UTC),
			"cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d",
		),
	}
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
