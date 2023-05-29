package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	junodb "github.com/forbole/juno/v5/database"
	"github.com/forbole/juno/v5/database/postgresql"
	juno "github.com/forbole/juno/v5/types"

	"github.com/desmos-labs/djuno/v2/x/authz"
	contracts "github.com/desmos-labs/djuno/v2/x/contracts/base"
	"github.com/desmos-labs/djuno/v2/x/contracts/tips"
	"github.com/desmos-labs/djuno/v2/x/feegrant"
	"github.com/desmos-labs/djuno/v2/x/notifications"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/profiles"
	profilesscore "github.com/desmos-labs/djuno/v2/x/profiles-score"
	"github.com/desmos-labs/djuno/v2/x/reactions"
	"github.com/desmos-labs/djuno/v2/x/relationships"
	"github.com/desmos-labs/djuno/v2/x/reports"
	"github.com/desmos-labs/djuno/v2/x/subspaces"
)

type Database interface {
	junodb.Database

	authz.Database
	contracts.Database
	tips.Database
	feegrant.Database
	notifications.Database
	posts.Database
	profiles.Database
	profilesscore.Database
	reactions.Database
	relationships.Database
	reports.Database
	subspaces.Database
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ Database = &Db{}
)

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type Db struct {
	cdc codec.Codec
	*postgresql.Database
}

// Builder allows to create a new Db instance implementing the database.Builder type
func Builder(ctx *junodb.Context) (junodb.Database, error) {
	database, err := postgresql.Builder(ctx)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid database type")
	}

	return &Db{
		cdc:      ctx.EncodingConfig.Codec,
		Database: psqlDb,
	}, nil
}

// Cast casts the given database to be a *Db
func Cast(database junodb.Database) Database {
	desmosDb, ok := (database).(Database)
	if !ok {
		panic(fmt.Errorf("database is not a DJuno database instance"))
	}
	return desmosDb
}

// --------------------------------------------------------------------------------------------------------------------

// SaveTx overrides postgresql.Database to perform a no-op
func (db *Db) SaveTx(_ *juno.Tx) error {
	return nil
}

// HasValidator overrides postgresql.Database to perform a no-op
func (db *Db) HasValidator(_ string) (bool, error) {
	return true, nil
}

// SaveValidators overrides postgresql.Database to perform a no-op
func (db *Db) SaveValidators(_ []*juno.Validator) error {
	return nil
}

// SaveCommitSignatures overrides postgresql.Database to perform a no-op
func (db *Db) SaveCommitSignatures(_ []*juno.CommitSig) error {
	return nil
}
