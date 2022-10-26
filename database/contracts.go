package database

import (
	"github.com/desmos-labs/djuno/v2/types"
)

// SaveContract stores the given contract data into the database
func (db *Db) SaveContract(contract types.Contract) error {
	stmt := `
INSERT INTO contract (address, type, height)
VALUES ($1, $2, $3)
ON CONFLICT (address) DO UPDATE 
    SET address = excluded.address, 
        type = excluded.type, 
        height = excluded.height
WHERE contract.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, contract.Address, contract.Type, contract.Height)
	return err
}

// SaveContractConfig stores the given contract config inside the database
func (db *Db) SaveContractConfig(config types.ContractConfig) error {
	stmt := `
INSERT INTO contract_config (contract_address, config, height)
VALUES ($1, $2, $3)
ON CONFLICT (contract_address) DO UPDATE
	SET contract_address = excluded.contract_address,
	    config = excluded.config,
	    height = excluded.height
WHERE contract_config.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, config.Address, string(config.ConfigJSONBz), config.Height)
	return err
}
