package notifications

import (
	"github.com/forbole/juno/v2/modules"
	junocfg "github.com/forbole/juno/v2/types/config"

	"github.com/desmos-labs/djuno/database"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
	_ modules.TransactionModule          = &Module{}
	_ modules.MessageModule              = &Module{}
)

// Module represents the module that will send users the notifications when something happens
type Module struct {
	cfg *Config
	db  *database.Db
}

// NewModule returns a new Module instance
func NewModule(cfg junocfg.Config, db *database.Db) *Module {
	notificationsCfg, err := ParseConfig(cfg.GetBytes())
	if err != nil {
		panic(err)
	}

	return &Module{
		cfg: notificationsCfg,
		db:  db,
	}
}

// Name implements modules.Module
func (m Module) Name() string {
	return "notifications"
}
