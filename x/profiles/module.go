package profiles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v3/node"
	"google.golang.org/grpc"

	"github.com/desmos-labs/djuno/v2/database"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	node           node.Node
	profilesClient profilestypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(node node.Node, grpcConnection *grpc.ClientConn, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		node:           node,
		profilesClient: profilestypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}
