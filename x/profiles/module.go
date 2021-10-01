package profiles

import (
	"encoding/json"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	"github.com/desmos-labs/juno/modules/messages"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/djuno/database"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module        = &Module{}
	_ modules.GenesisModule = &Module{}
	_ modules.MessageModule = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	encodingConfig *params.EncodingConfig
	db             *database.Db
	profilesClient profilestypes.QueryClient
	getAccounts    messages.MessageAddressesParser
}

// NewModule allows to build a new Module instance
func NewModule(
	getAccounts messages.MessageAddressesParser, profilesClient profilestypes.QueryClient,
	encodingConfig *params.EncodingConfig, db *database.Db,
) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		db:             db,
		getAccounts:    getAccounts,
		profilesClient: profilesClient,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles"
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.profilesClient, m.encodingConfig.Marshaler, m.db)
}
