package reports

import (
	"context"
	"fmt"

	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"
	"github.com/forbole/juno/v3/node/remote"
)

// getStandardReason returns the standard reason having the given id at the given height
func (m *Module) getStandardReason(height int64, id uint32) (reportstypes.StandardReason, error) {
	res, err := m.reportsClient.Params(
		remote.GetHeightRequestContext(context.Background(), height),
		&reportstypes.QueryParamsRequest{},
	)
	if err != nil {
		return reportstypes.StandardReason{}, err
	}

	for _, reason := range res.Params.StandardReasons {
		if reason.ID == id {
			return reason, nil
		}
	}

	return reportstypes.StandardReason{}, fmt.Errorf("standard reason with id %d not found", id)
}
