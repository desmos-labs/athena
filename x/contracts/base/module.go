package contracts

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	"github.com/desmos-labs/djuno/v2/database"
)

// Module represents a generic smart contract module that can be extended for custom contracts handling
type Module struct {
	wasmClient wasmtypes.QueryClient
	db         *database.Db
}

// NewModule returns a new Module instance
func NewModule(wasmClient wasmtypes.QueryClient, db *database.Db) *Module {
	return &Module{
		wasmClient: wasmClient,
		db:         db,
	}
}
