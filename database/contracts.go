package database

import (
	"fmt"

	"github.com/desmos-labs/djuno/v2/types"
	contracts "github.com/desmos-labs/djuno/v2/x/contracts/base"
)

var (
	_ contracts.Database = &Db{}
)

// SaveContract stores the given contract data into the database
func (db *Db) SaveContract(contract types.Contract) error {
	stmt := `
INSERT INTO contract (address, type, config, height)
VALUES ($1, $2, $3, $4)
ON CONFLICT (address) DO UPDATE 
    SET address = excluded.address, 
        type = excluded.type,
        config = excluded.config,
        height = excluded.height
WHERE contract.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, contract.Address, contract.Type, string(contract.ConfigBz), contract.Height)
	return err
}

type contractRow struct {
	Address  string `db:"address"`
	Type     string `db:"type"`
	ConfigBz []byte `db:"config"`
	Height   int64  `db:"height"`
}

// GetContract returns the stored data of the contract with the given address
func (db *Db) GetContract(address string) (*types.Contract, error) {
	var rows []contractRow
	stmt := `SELECT * FROM contract WHERE address = $1`
	err := db.SQL.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	if len(rows) > 1 {
		return nil, fmt.Errorf("multiple contracts found for address %s", err)
	}

	row := rows[0]
	contract := types.NewContract(row.Address, row.Type, row.ConfigBz, row.Height)
	return &contract, nil
}
