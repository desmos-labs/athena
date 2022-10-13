package authz

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Info().Str("module", "authz").Msg("setting up periodic tasks")

	// Delete expired grants every 5 minutes
	if _, err := scheduler.Every(5).Minutes().StartImmediately().Do(m.deleteExpiredGrants); err != nil {
		return fmt.Errorf("error while scheduling authz peridic operation: %s", err)
	}

	return nil
}

// deleteExpiredGrants deletes the expired grants from the database
func (m *Module) deleteExpiredGrants() {
	err := m.db.DeleteExpiredGrants(time.Now())
	if err != nil {
		log.Error().Str("module", "authz").Err(err).Msg("error while deleting expired grants")
	}
}
