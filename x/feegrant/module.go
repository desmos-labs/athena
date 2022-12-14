package feegrant

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
)

var (
	_ modules.Module        = &Module{}
	_ modules.MessageModule = &Module{}
)

type Module struct {
	node node.Node
	cdc  codec.Codec
	db   Database
}

// NewModule returns a new Module instance
func NewModule(node node.Node, cdc codec.Codec, db Database) *Module {
	return &Module{
		node: node,
		cdc:  cdc,
		db:   db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "feegrant"
}
