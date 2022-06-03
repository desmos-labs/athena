package fees

import (
	"context"
	"fmt"

	feestypes "github.com/desmos-labs/desmos/v3/x/fees/types"
	"github.com/forbole/juno/v3/node/remote"

	"github.com/desmos-labs/djuno/v2/types"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// RegisterPeriodicOperations implements modules.PeriodicOperationsModule
func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Info().Str("module", "fees").Msg("setting up periodic tasks")

	// Update the params every 30 mins
	if _, err := scheduler.Every(30).Minutes().StartImmediately().Do(m.updateParams); err != nil {
		return fmt.Errorf("error while scheduling fees peridic operation: %s", err)
	}

	return nil
}

// updateParams allows to update the fees params by fetching them from the chain
func (m *Module) updateParams() error {
	height, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	res, err := m.feesClient.Params(
		remote.GetHeightRequestContext(context.Background(), height),
		&feestypes.QueryParamsRequest{},
	)
	if err != nil {
		return fmt.Errorf("error while getting params: %s", err)
	}

	err = m.db.SaveFeesParams(types.NewFeesParams(res.Params, height))
	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}
