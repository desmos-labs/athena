package database

import (
	"fmt"

	"github.com/forbole/juno/v3/database"
	"github.com/forbole/juno/v3/database/postgresql"
	juno "github.com/forbole/juno/v3/types"
)

// Db represents a PostgreSQL database with expanded features.
// so that it can properly store posts and other Desmos-related data.
type Db struct {
	*postgresql.Database
}

// Cast casts the given database to be a *Db
func Cast(database database.Database) *Db {
	desmosDb, ok := (database).(*Db)
	if !ok {
		panic(fmt.Errorf("database is not a DesmosDB instance"))
	}
	return desmosDb
}

// Builder allows to create a new Db instance implementing the database.Builder type
func Builder(ctx *database.Context) (database.Database, error) {
	database, err := postgresql.Builder(ctx)
	if err != nil {
		return nil, err
	}

	psqlDb, ok := (database).(*postgresql.Database)
	if !ok {
		return nil, fmt.Errorf("invalid database type")
	}

	return &Db{
		Database: psqlDb,
	}, nil
}

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

// SaveMessage overrides postgresql.Database to perform a no-op
func (db *Db) SaveMessage(_ *juno.Message) error {
	return nil
}
