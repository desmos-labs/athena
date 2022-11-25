package profilesscore

import (
	"github.com/forbole/juno/v4/modules"
)

var (
	_ modules.Module                   = &Module{}
	_ modules.PeriodicOperationsModule = &Module{}
)

type Module struct {
	scorers Scorers
	db      Database
}

// NewModule returns a new Module instance
func NewModule(scorers Scorers, db Database) *Module {
	return &Module{
		scorers: scorers,
		db:      db,
	}
}

// Name implements modules.Module
func (m *Module) Name() string {
	return "profiles:score"
}

func (m *Module) GetScorers() Scorers {
	var scorers Scorers
	for _, scorer := range m.scorers {
		if scorer != nil {
			scorers = append(scorers, scorer)
		}
	}
	return scorers
}
