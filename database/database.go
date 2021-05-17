package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	juno "github.com/desmos-labs/juno/types"
	"github.com/jmoiron/sqlx"
)

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type Db struct {
	*postgresql.Database
	Sqlx *sqlx.DB
}

// Cast casts the given database to be a *Db
func Cast(database db.Database) *Db {
	desmosDb, ok := (database).(*Db)
	if !ok {
		panic(fmt.Errorf("database is not a DesmosDB instance"))
	}
	return desmosDb
}

// Builder allows to create a new Db instance implementing the database.Builder type
func Builder(cfg juno.Config, encodingConfig *params.EncodingConfig) (db.Database, error) {
	database, err := postgresql.Builder(cfg.GetDatabaseConfig(), encodingConfig)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid database type")
	}

	return &Db{
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}, nil
}
