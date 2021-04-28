package reports

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"
)

var _ modules.Module = &Module{}
var _ modules.MessageModule = &Module{}

// Module represents the x/reports module handler
type Module struct {
	db *database.DesmosDb
}

// NewModule returns a new Module instance
func NewModule(db *database.DesmosDb) *Module {
	return &Module{
		db: db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "reports"
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	return HandleMsg(tx, msg, m.db)
}
