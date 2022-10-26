package database

import (
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveContract stores the given contract data into the database
func (db *Db) SaveContract(contract types.Contract) error {
	stmt := `
INSERT INTO contract (address, type, config, height)
VALUES ($1, $2, $4, $3)
ON CONFLICT (address) DO UPDATE 
    SET address = excluded.address, 
        type = excluded.type,
        config = excluded.config,
        height = excluded.height
WHERE contract.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, contract.Address, contract.Type, string(contract.ConfigBz), contract.Height)
	return err
}
