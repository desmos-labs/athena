package posts

import (
	"github.com/cosmos/cosmos-sdk/codec"
	poststypes "github.com/desmos-labs/desmos/v4/x/posts/types"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/node"
	"google.golang.org/grpc"

	"github.com/desmos-labs/djuno/v2/database"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.BlockModule              = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/fees module handler
type Module struct {
	cdc    codec.Codec
	db     *database.Db
	node   node.Node
	client poststypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(node node.Node, grpcConnection *grpc.ClientConn, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		node:   node,
		client: poststypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "posts"
}
