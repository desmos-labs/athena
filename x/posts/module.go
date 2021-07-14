package posts

import (
	"encoding/json"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"

	"github.com/desmos-labs/djuno/database"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var _ modules.Module = &Module{}
var _ modules.GenesisModule = &Module{}
var _ modules.MessageModule = &Module{}

// Module represents the x/posts module handler
type Module struct {
	encodingConfig *params.EncodingConfig
	db             *database.Db
	profilesClient profilestypes.QueryClient
}

// NewModule allows to build a new Module instance
func NewModule(profilesClient profilestypes.QueryClient, encodingConfig *params.EncodingConfig, db *database.Db) *Module {
	return &Module{
		encodingConfig: encodingConfig,
		db:             db,
		profilesClient: profilesClient,
	}
}

// Name implements Module
func (m *Module) Name() string {
	return "posts"
}

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	return HandleGenesis(doc, appState, m.encodingConfig.Marshaler, m.db)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return MsgHandler(tx, index, msg, m.profilesClient, m.encodingConfig.Marshaler, m.db)
}
