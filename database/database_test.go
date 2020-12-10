package database_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	poststypes "github.com/desmos-labs/desmos/x/posts/types"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

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
	// Setup the test data
	creator, err := sdk.AccAddressFromBech32("cosmos1qpzgtwec63yhxz9hesj8ve0j3ytzhhqaqxrc5d")
	suite.Require().NoError(err)

	created, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)

	suite.testData = TestData{
		post: poststypes.NewPost(
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
		),
	}
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
