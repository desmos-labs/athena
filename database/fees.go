package database

import (
	"encoding/json"
	"fmt"

	"github.com/desmos-labs/djuno/v2/types"
)

// SaveFeesParams allows to store the given fees params
func (db *Db) SaveFeesParams(params types.FeesParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling fees params: %s", err)
	}

	stmt := `
INSERT INTO fees_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE fees_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing fees params: %s", err)
	}

	return nil
}
