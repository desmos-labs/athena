package notifications

import (
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"

	"github.com/desmos-labs/djuno/types"
)

var _ modules.Module = &Module{}
var _ modules.TransactionModule = &Module{}

// Module represents the module that will send users the notifications when something happens
type Module struct{}

// NewModule returns a new Module instance
func NewModule(cfg *types.Config) *Module {
	if cfg.Notifications.Enable {
		err := setupNotifications(cfg.Notifications)
		if err != nil {
			panic(err)
		}
	}

	return &Module{}
}

// Name implements modules.Module
func (m Module) Name() string {
	return "notifications"
}

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	return TxHandler(tx)
}
