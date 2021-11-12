package profiles

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "profiles").Msg("setting up periodic tasks")

	// Update the params every 30 mins
	if _, err := scheduler.Every(30).Minutes().Do(m.updateParams); err != nil {
		return fmt.Errorf("error while scheduling profiles peridic operation: %s", err)
	}

	return nil
}
