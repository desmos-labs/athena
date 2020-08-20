package database

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/djuno/cmd/djuno/flags"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/db/postgresql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// DesmosDb represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type DesmosDb struct {
	postgresql.Database
	Sqlx *sqlx.DB
}

// Builder allows to create a new DesmosDb instance implementing the database.Builder type
func Builder(cfg config.Config, codec *codec.Codec) (*db.Database, error) {
	psqlConfg, ok := cfg.DatabaseConfig.Config.(*config.PostgreSQLConfig)
	if !ok {
		// TODO: Support MongoDB
		return nil, fmt.Errorf("mongodb configuration is not supported on Djuno")
	}

	// Override config values if needed
	overrideValues(psqlConfg)

	database, err := postgresql.Builder(*psqlConfg, codec)
	if err != nil {
		return nil, err
	}

	psqlDb, _ := (*database).(postgresql.Database)
	var desmosDb db.Database = DesmosDb{
		Database: psqlDb,
		Sqlx:     sqlx.NewDb(psqlDb.Sql, "postgresql"),
	}

	return &desmosDb, nil
}

// overrideValues takes the given cfg and replaces any default value with the ones specified using the command flags
func overrideValues(cfg *config.PostgreSQLConfig) {
	overridingDbHost := viper.GetString(flags.FlagDBHost)
	if len(overridingDbHost) != 0 {
		cfg.Host = overridingDbHost
	}

	overridingDbPort := viper.GetUint64(flags.FlagDBPort)
	if overridingDbPort > 0 {
		cfg.Port = overridingDbPort
	}

	overridingDbUser := viper.GetString(flags.FlagDBUser)
	if len(overridingDbUser) != 0 {
		cfg.User = overridingDbUser
	}

	overridingDbPassword := viper.GetString(flags.FlagDBPassword)
	if len(overridingDbPassword) != 0 {
		cfg.Password = overridingDbPassword
	}

	overridingDbName := viper.GetString(flags.FlagDBName)
	if len(overridingDbName) != 0 {
		cfg.Name = overridingDbName
	}
}
