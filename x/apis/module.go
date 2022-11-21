package apis

import (
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/types/config"
)

var (
	_ modules.Module                     = &Module{}
	_ modules.AdditionalOperationsModule = &Module{}
)

type Module struct {
	cfg       *Config
	registrar Registrar
}

func NewModule(junoCfg config.Config, registrar Registrar) *Module {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}
	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	return &Module{
		cfg:       cfg,
		registrar: registrar,
	}
}

func (m Module) Name() string {
	return "apis"
}
