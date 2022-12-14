package apis

import (
	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

// Module represents the module allowing to register custom API endpoints
type Module struct {
	ctx          Context
	cfg          *Config
	registrar    Registrar
	configurator Configurator
}

func NewModule(ctx Context) *Module {
	cfgBz, err := ctx.JunoConfig.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	return &Module{
		ctx: ctx,
		cfg: cfg,
	}
}

// WithRegistrar allows setting the APIs registrar to be used
func (m *Module) WithRegistrar(registrar Registrar) *Module {
	m.registrar = registrar
	return m
}

// WithConfigurator allows setting the configurator to be used
func (m *Module) WithConfigurator(configurator Configurator) *Module {
	m.configurator = configurator
	return m
}

func (m *Module) Name() string {
	return "apis"
}
