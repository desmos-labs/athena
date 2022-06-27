package reports

import (
	"context"
	"github.com/desmos-labs/djuno/v2/types"

	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"
	"github.com/forbole/juno/v3/node/remote"
)

// updateReport updates the stored data for the given report at the specified height
func (m *Module) updateReport(height int64, subspaceID uint64, reportID uint64) error {
	// Get the report
	res, err := m.reportsClient.Report(
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
	res, err := m.reportsClient.Reason(
		remote.GetHeightRequestContext(context.Background(), height),
		&reportstypes.QueryReasonRequest{SubspaceId: subspaceID, ReasonId: reasonID},
	)
	if err != nil {
		return err
	}

	// Save the reason
	return m.db.SaveReason(types.NewReason(res.Reason, height))
}
