package database

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
)

// DesmosDb represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type DesmosDb struct {
	*postgresql.Database
	Sqlx *sqlx.DB
}

// Cast casts the given database to be a *DesmosDb
func Cast(database db.Database) *DesmosDb {
	desmosDb, ok := (database).(*DesmosDb)
	if !ok {
		panic(fmt.Errorf("database is not a DesmosDB instance"))
	}
	return desmosDb
}

// Builder allows to create a new DesmosDb instance implementing the database.Builder type
func Builder(cfg *config.Config, encodingConfig *params.EncodingConfig) (db.Database, error) {
	database, err := postgresql.Builder(cfg.Database, encodingConfig)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid database type")
	}

	return &DesmosDb{
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}, nil
}
