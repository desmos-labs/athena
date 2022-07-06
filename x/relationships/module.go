package relationships

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"google.golang.org/grpc"

	relationshipstypes "github.com/desmos-labs/desmos/v4/x/relationships/types"

	"github.com/desmos-labs/djuno/v2/database"

	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	cdc            codec.Codec
	db             *database.Db
	profilesModule ProfilesModule
	client         relationshipstypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(profilesModule ProfilesModule, grpcConnection *grpc.ClientConn, cdc codec.Codec, db *database.Db) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		profilesModule: profilesModule,
		client:         relationshipstypes.NewQueryClient(grpcConnection),
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "relationships"
}
