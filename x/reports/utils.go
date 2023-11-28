package reports

import (
	"context"
	"fmt"

	"github.com/desmos-labs/athena/types"

	"github.com/forbole/juno/v5/node/remote"

	reportstypes "github.com/desmos-labs/desmos/v6/x/reports/types"
)

// updateReport updates the stored data for the given report at the specified height
func (m *Module) updateReport(height int64, subspaceID uint64, reportID uint64) error {
	// Get the report
	res, err := m.client.Report(
		remote.GetHeightRequestContext(context.Background(), height),
		&reportstypes.QueryReportRequest{SubspaceId: subspaceID, ReportId: reportID},
	)
	if err != nil {
		return err
	}

	// Save the report
	return m.db.SaveReport(types.NewReport(res.Report, height))
}

// updateReason updates the stored data for the given reason at the specified height
func (m *Module) updateReason(height int64, subspaceID uint64, reasonID uint32) error {
	// Get the reason
	res, err := m.client.Reason(
		remote.GetHeightRequestContext(context.Background(), height),
		&reportstypes.QueryReasonRequest{SubspaceId: subspaceID, ReasonId: reasonID},
	)
	if err != nil {
		return err
	}

	// Save the reason
	return m.db.SaveReason(types.NewReason(res.Reason, height))
}

// updateParams updates the stored params for the given height
func (m *Module) updateParams() error {
	height, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block height: %s", err)
	}

	// Get the params
	res, err := m.client.Params(
		remote.GetHeightRequestContext(context.Background(), height),
		&reportstypes.QueryParamsRequest{},
	)
	if err != nil {
		return err
	}

	// Save the params
	return m.db.SaveReportsParams(types.NewReportsParams(res.Params, height))
}
