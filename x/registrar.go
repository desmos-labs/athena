package x

import (
	"fmt"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/modules"
	juno "github.com/desmos-labs/juno/types"

	"github.com/desmos-labs/djuno/types"
	"github.com/desmos-labs/djuno/x/common"

	"github.com/desmos-labs/djuno/database"
	"github.com/desmos-labs/djuno/x/notifications"
	"github.com/desmos-labs/djuno/x/posts"
	"github.com/desmos-labs/djuno/x/profiles"
)

// ModulesRegistrar represents the modules.Registrar that allows to register all custom DJuno modules
type ModulesRegistrar struct {
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar() *ModulesRegistrar {
	return &ModulesRegistrar{}
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(
	cfg juno.Config, encodingConfig *params.EncodingConfig, _ *sdk.Config, db db.Database, _ *client.Proxy,
) modules.Modules {
	desmosDb := database.Cast(db)

	djunoCfg, ok := cfg.(*types.Config)
	if !ok {
		panic(fmt.Errorf("invalid configuration type: %T", cfg))
	}

	grpcConnection := client.MustCreateGrpcConnection(cfg)

	return []modules.Module{
		notifications.NewModule(djunoCfg.Notifications, desmosDb),
		posts.NewModule(encodingConfig, desmosDb),
		profiles.NewModule(common.MessagesParser, profilestypes.NewQueryClient(grpcConnection), encodingConfig, desmosDb),
	}
}
