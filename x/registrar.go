package x

import (
	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/bank"
	"github.com/desmos-labs/djuno/x/notifications"
	"github.com/desmos-labs/djuno/x/posts"
	"github.com/desmos-labs/djuno/x/profiles"
	"github.com/desmos-labs/djuno/x/reports"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
)

// ModulesRegistrar represents the modules.Registrar that allows to register all custom BDJuno modules
type ModulesRegistrar struct {
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar() *ModulesRegistrar {
	return &ModulesRegistrar{}
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(
	cfg *config.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, cp *client.Proxy,
) modules.Modules {
	desmosDb := database.Cast(db)
	return []modules.Module{
		bank.NewModule(),
		notifications.NewModule(),
		posts.NewModule(encodingConfig, desmosDb),
		profiles.NewModule(encodingConfig, desmosDb),
		reports.NewModule(desmosDb),
	}
}
