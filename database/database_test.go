package database_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	junodb "github.com/forbole/juno/v3/database"
	junodbcfg "github.com/forbole/juno/v3/database/config"
	"github.com/forbole/juno/v3/logging"

	desmosapp "github.com/desmos-labs/desmos/v4/app"
	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/djuno/v2/database"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database *database.Db
}

func (suite *DbTestSuite) SetupTest() {
	// Build the database
	encodingConfig := desmosapp.MakeTestEncodingConfig()
	databaseConfig := junodbcfg.NewDatabaseConfig(
		"djuno",
		"localhost",
		6432,
		"djuno",
		"password",
		"",
		"public",
		10,
		10,
		0,
		0,
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

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
