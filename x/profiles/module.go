package profiles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/forbole/juno/v5/node"
	"google.golang.org/grpc"

	profilestypes "github.com/desmos-labs/desmos/v6/x/profiles/types"

	"github.com/forbole/juno/v5/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
	_ modules.GenesisModule            = &Module{}
	_ modules.MessageModule            = &Module{}
	_ modules.AuthzMessageModule       = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	cdc            codec.Codec
	db             Database
	node           node.Node
	authClient     authtypes.QueryClient
	profilesClient profilestypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(node node.Node, grpcConnection *grpc.ClientConn, cdc codec.Codec, db Database) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		node:           node,
		authClient:     authtypes.NewQueryClient(grpcConnection),
		profilesClient: profilestypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}
