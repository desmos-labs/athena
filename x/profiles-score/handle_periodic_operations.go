package profilesscore

import (
	"fmt"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Info().Str("module", "profiles score").Msg("setting up periodic tasks")

	// Update the scores every day
	if _, err := scheduler.Every(1).Days().StartImmediately().Do(m.updateApplicationLinkScores); err != nil {
		return fmt.Errorf("error while scheduling profiles score peridic operation: %s", err)
	}

	return nil
}

// updateApplicationLinkScores updates the score for each of the stored application links
func (m *Module) updateApplicationLinkScores() {
	err := m.RefreshApplicationLinksScores()
	if err != nil {
		log.Error().Err(err).Msg("error while refreshing applications links scores")
	}
}
