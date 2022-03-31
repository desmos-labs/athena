package x

import (
	"fmt"

	"github.com/desmos-labs/djuno/v2/x/relationships"

	"github.com/forbole/juno/v3/node/builder"

	"github.com/forbole/juno/v3/node/remote"

	"github.com/forbole/juno/v3/modules/registrar"

	"github.com/forbole/juno/v3/modules"

	"github.com/desmos-labs/djuno/v2/database"
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
	cdc := ctx.EncodingConfig.Marshaler
	desmosDb := database.Cast(ctx.Database)

	remoteCfg, ok := ctx.JunoConfig.Node.Details.(*remote.Details)
	if !ok {
		panic(fmt.Errorf("cannot run DJuno on local node"))
	}

	node, err := builder.BuildNode(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(fmt.Errorf("cannot build node: %s", err))
	}

	grpcConnection := remote.MustCreateGrpcConnection(remoteCfg.GRPC)
	profilesModule := profiles.NewModule(node, grpcConnection, cdc, desmosDb)
	relationshipsModule := relationships.NewModule(profilesModule, grpcConnection, cdc, desmosDb)

	return []modules.Module{
		profilesModule,
		relationshipsModule,
	}
}
