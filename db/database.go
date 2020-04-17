package db

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
)

// DesmosDb represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type DesmosDb struct {
	postgresql.Database
	sqlx *sqlx.DB
}

// Builder allows to create a new DesmosDb instance implementing the database.Builder type
func Builder(cfg config.Config, codec *codec.Codec) (*db.Database, error) {
	psqlConfg, ok := cfg.DatabaseConfig.Config.(*config.PostgreSQLConfig)
	if !ok {
		// TODO: Support MongoDB
		return nil, fmt.Errorf("mongodb configuration is not supported on Djuno")
	}

	database, err := postgresql.Builder(*psqlConfg, codec)
	if err != nil {
		return nil, err
	}

	psqlDb, _ := (*database).(postgresql.Database)
	var desmosDb db.Database = DesmosDb{
		Database: psqlDb,
		sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}

	return &desmosDb, nil
}
