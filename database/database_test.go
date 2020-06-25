package database_test

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/djuno/database"
	jconfig "github.com/desmos-labs/juno/config"
	"github.com/stretchr/testify/suite"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database database.DesmosDb
}

func (suite *DbTestSuite) SetupTest() {
	// Create the codec
	codec := desmosapp.MakeCodec()

	// Build the database
	config := jconfig.Config{
		DatabaseConfig: jconfig.DatabaseConfig{
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

	desmosDb, ok := (*db).(database.DesmosDb)
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

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
