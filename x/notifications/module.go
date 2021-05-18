package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"

	"github.com/desmos-labs/djuno/database"

	"github.com/desmos-labs/djuno/types"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
	_ modules.TransactionModule          = &Module{}
	_ modules.MessageModule              = &Module{}
)

// Module represents the module that will send users the notifications when something happens
type Module struct {
	cfg *types.NotificationsConfig
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg *types.NotificationsConfig, db *database.Db) *Module {
	return &Module{
		cfg: cfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m Module) Name() string {
	return "notifications"
}

// RunAdditionalOperations implements modules.AdditionalOperationsModule
func (m Module) RunAdditionalOperations() error {
	return setupNotifications(m.cfg)
}

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	return TxHandler(tx)
}

// HandleMsg implements modules.MessageModule
func (m Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	return MsgHandler(tx, index, msg, m.db)
}
