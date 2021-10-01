package x

import (
	"fmt"

	"github.com/desmos-labs/juno/modules/registrar"

	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	"github.com/desmos-labs/juno/client"
	"github.com/desmos-labs/juno/modules"

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
func (r *ModulesRegistrar) BuildModules(ctx registrar.Context) modules.Modules {
	desmosDb := database.Cast(ctx.Database)

	djunoCfg, ok := ctx.ParsingConfig.(*types.Config)
	if !ok {
		panic(fmt.Errorf("invalid configuration type: %T", ctx.ParsingConfig))
	}

	grpcConnection := client.MustCreateGrpcConnection(ctx.ParsingConfig)
	profilesClient := profilestypes.NewQueryClient(grpcConnection)

	return []modules.Module{
		notifications.NewModule(djunoCfg.Notifications, desmosDb),
		posts.NewModule(profilesClient, ctx.EncodingConfig, desmosDb),
		profiles.NewModule(common.MessagesParser, profilesClient, ctx.EncodingConfig, desmosDb),
	}
}
