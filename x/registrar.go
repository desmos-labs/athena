package x

import (
	"fmt"

	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/juno/v2/modules/registrar"

	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	"github.com/forbole/juno/v2/modules"

	"github.com/desmos-labs/djuno/v2/x/common"

	"github.com/desmos-labs/djuno/v2/database"
	"github.com/desmos-labs/djuno/v2/x/notifications"
	"github.com/desmos-labs/djuno/v2/x/posts"
	"github.com/desmos-labs/djuno/v2/x/profiles"
)

// ModulesRegistrar represents the modules.Registrar that allows to register all custom DJuno modules
type ModulesRegistrar struct {
}

// NewModulesRegistrar allows to build a new ModulesRegistrar instance
func NewModulesRegistrar() *ModulesRegistrar {
	return &ModulesRegistrar{}
}

// BuildModules implements modules.Registrar
func (r *ModulesRegistrar) BuildModules(ctx registrar.Context) modules.Modules {
	desmosDb := database.Cast(ctx.Database)

	remoteCfg, ok := ctx.JunoConfig.Node.Details.(*remote.Details)
	if !ok {
		panic(fmt.Errorf("cannot run DJuno on local node"))
	}

	grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
	profilesClient := profilestypes.NewQueryClient(grpcConnection)

	profilesModule := profiles.NewModule(common.MessagesParser, profilesClient, ctx.EncodingConfig.Marshaler, desmosDb)

	return []modules.Module{
		profilesModule,
		notifications.NewModule(ctx.JunoConfig, desmosDb),
		posts.NewModule(ctx.EncodingConfig.Marshaler, desmosDb, profilesModule),
	}
}
