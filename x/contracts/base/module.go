package contracts

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// Module represents a generic smart contract module that can be extended for custom contracts handling
type Module struct {
	wasmClient wasmtypes.QueryClient
	db         Database
}

// NewModule returns a new Module instance
func NewModule(wasmClient wasmtypes.QueryClient, db Database) *Module {
	return &Module{
		wasmClient: wasmClient,
		db:         db,
	}
}
