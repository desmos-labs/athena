package fees

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/node"
	"google.golang.org/grpc"

	feestypes "github.com/desmos-labs/desmos/v4/x/fees/types"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/fees module handler
type Module struct {
	cdc        codec.Codec
	db         Database
	node       node.Node
	feesClient feestypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(node node.Node, grpcConnection *grpc.ClientConn, cdc codec.Codec, db Database) *Module {
	return &Module{
		cdc:        cdc,
		db:         db,
		node:       node,
		feesClient: feestypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "fees"
}
