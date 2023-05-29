package reports

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/node"
	"google.golang.org/grpc"

	reportstypes "github.com/desmos-labs/desmos/v5/x/reports/types"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

// Module represents the x/fees module handler
type Module struct {
	cdc    codec.Codec
	db     Database
	node   node.Node
	client reportstypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(node node.Node, grpcConnection *grpc.ClientConn, cdc codec.Codec, db Database) *Module {
	return &Module{
		cdc:    cdc,
		db:     db,
		node:   node,
		client: reportstypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "reports"
}
