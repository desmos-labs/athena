package apis

import (
	"github.com/forbole/juno/v4/modules"
	"github.com/forbole/juno/v4/modules/registrar"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

// Module represnets the module allowing to register custom API endpoints
type Module struct {
	ctx       registrar.Context
	cfg       *Config
	registrar Registrar
}

func NewModule(ctx registrar.Context, registrar Registrar) *Module {
	cfgBz, err := ctx.JunoConfig.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	return &Module{
		ctx:       ctx,
		cfg:       cfg,
		registrar: registrar,
	}
}

func (m Module) Name() string {
	return "apis"
}
