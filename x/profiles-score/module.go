package profilesscore

import (
	"github.com/desmos-labs/djuno/v2/x/profiles-score/scorers"
	"github.com/forbole/juno/v3/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

type Module struct {
	scorers []scorers.Scorer
}

func NewModule() *Module {
	return &Module{}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles:score"
}
