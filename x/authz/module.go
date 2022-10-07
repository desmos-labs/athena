package authz

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/node"

	"github.com/desmos-labs/djuno/v2/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

type Module struct {
	cdc  codec.Codec
	node node.Node
	db   *database.Db
}

// NewModule returns a new Module instance
func NewModule(node node.Node, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		node: node,
		cdc:  cdc,
		db:   db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "authz"
}
