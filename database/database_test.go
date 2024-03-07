package database_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/desmos-labs/desmos/v7/app"
	profilestypes "github.com/desmos-labs/desmos/v7/x/profiles/types"
	junodb "github.com/forbole/juno/v5/database"
	junodbcfg "github.com/forbole/juno/v5/database/config"
	"github.com/forbole/juno/v5/logging"
	"github.com/forbole/juno/v5/types/params"

	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/athena/v2/database"

	_ "github.com/proullon/ramsql/driver"
)

type DbTestSuite struct {
	suite.Suite

	database *database.Db
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}

func (suite *DbTestSuite) SetupSuite() {
	// Build the database
	encodingConfig := app.MakeEncodingConfig()
	databaseConfig := junodbcfg.DefaultDatabaseConfig().
		WithURL("postgres://athena:password@localhost:6432/athena?sslmode=disable&search_path=public")

	db, err := database.Builder(junodb.NewContext(databaseConfig, params.EncodingConfig(encodingConfig), logging.DefaultLogger()))
	suite.Require().NoError(err)

	desmosDb, ok := (db).(*database.Db)
	suite.Require().True(ok)

	// Delete the public schema
	_, err = desmosDb.SQL.Exec(fmt.Sprintf(`DROP SCHEMA %s CASCADE;`, databaseConfig.GetSchema()))
	suite.Require().NoError(err)

	// Create the schema
	_, err = desmosDb.SQL.Exec(fmt.Sprintf(`CREATE SCHEMA %s;`, databaseConfig.GetSchema()))
	suite.Require().NoError(err)

	dirPath := "schema"
	dir, err := os.ReadDir(dirPath)
	for _, fileInfo := range dir {
		if !strings.HasSuffix(fileInfo.Name(), ".sql") {
			continue
		}

		file, err := os.ReadFile(filepath.Join(dirPath, fileInfo.Name()))
		suite.Require().NoError(err)

		commentsRegExp := regexp.MustCompile(`/\*.*\*/`)
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			_, err := desmosDb.SQL.Exec(commentsRegExp.ReplaceAllString(request, ""))
			suite.Require().NoError(err)
		}
	}

	// Create the truncate function
	stmt := fmt.Sprintf(`
CREATE OR REPLACE FUNCTION truncate_tables(username IN VARCHAR) RETURNS void AS $$
DECLARE
    statements CURSOR FOR
        SELECT tablename FROM pg_tables
        WHERE tableowner = username AND schemaname = '%s';
BEGIN
    FOR stmt IN statements LOOP
        EXECUTE 'TRUNCATE TABLE ' || quote_ident(stmt.tablename) || ' CASCADE;';
    END LOOP;
END;
$$ LANGUAGE plpgsql;`, databaseConfig.GetSchema())
	_, err = desmosDb.SQL.Exec(stmt)
	suite.Require().NoError(err)

	suite.database = desmosDb
}

func (suite *DbTestSuite) SetupTest() {
	_, err := suite.database.SQL.Exec(`SELECT truncate_tables('athena')`)
	suite.Require().NoError(err)
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) buildProfile(address string) *profilestypes.Profile {
	addr, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)

	profile, err := profilestypes.NewProfile(
		"TestUser",
		"Test User",
		"This is a test user",
		profilestypes.NewPictures(
			"https://example.com/profile.png",
			"https://example.com/cover.png",
		),
		time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		authtypes.NewBaseAccountWithAddress(addr),
	)
	suite.Require().NoError(err)

	return profile
}
