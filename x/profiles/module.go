package profiles

import (
	"github.com/cosmos/cosmos-sdk/codec"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"
	"github.com/forbole/juno/v2/modules/messages"

	"github.com/desmos-labs/djuno/database"

	"github.com/forbole/juno/v2/modules"
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
	profilesClient profilestypes.QueryClient
	getAccounts    messages.MessageAddressesParser
}

// NewModule allows to build a new Module instance
func NewModule(
	getAccounts messages.MessageAddressesParser, profilesClient profilestypes.QueryClient, cdc codec.Codec, db *database.Db,
) *Module {
	return &Module{
		cdc:            cdc,
		db:             db,
		getAccounts:    getAccounts,
		profilesClient: profilesClient,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}
