package profiles

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/modules/messages"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/djuno/database"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ modules.Module            = &Module{}
	_ modules.GenesisModule     = &Module{}
	_ modules.TransactionModule = &Module{}
	_ modules.MessageModule     = &Module{}
)

// Module represents the x/profiles module handler
type Module struct {
	encodingConfig *params.EncodingConfig
	cdc            codec.Marshaler
	db             *database.Db
	getAccounts    messages.MessageAddressesParser
}

// NewModule allows to build a new Module instance
func NewModule(
	getAccounts messages.MessageAddressesParser, encodingConfig *params.EncodingConfig, db *database.Db,
) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		db:             db,
		getAccounts:    getAccounts,
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

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	panic("implement me")
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, index, msg, m.getAccounts, m.encodingConfig.Marshaler, m.db)
}
