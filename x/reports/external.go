package reports

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/forbole/juno/v5/node/remote"

	reportstypes "github.com/desmos-labs/desmos/v5/x/reports/types"

	"github.com/desmos-labs/djuno/v2/types"
)

// RefreshReportsData refreshes all the reports data for the given subspace
func (m *Module) RefreshReportsData(height int64, subspaceID uint64) error {
	reports, err := m.queryAllReports(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllReports(height, subspaceID)
	if err != nil {
		return err
	}

	for _, report := range reports {
		log.Info().Uint64("subspace", report.SubspaceID).Uint64("report", report.ID).Msg("refreshing report")

		err = m.db.SaveReport(report)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllReports queries all the reports for the given subspace from the node
func (m *Module) queryAllReports(height int64, subspaceID uint64) ([]types.Report, error) {
	var reports []types.Report

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Reports(
			remote.GetHeightRequestContext(context.Background(), height),
			&reportstypes.QueryReportsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, report := range res.Reports {
			reports = append(reports, types.NewReport(report, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return reports, nil
}

// --------------------------------------------------------------------------------------------------------------------

// RefreshReasonsData refreshes all the reasons data for the given subspace
func (m *Module) RefreshReasonsData(height int64, subspaceID uint64) error {
	reasons, err := m.queryAllReasons(height, subspaceID)
	if err != nil {
		return err
	}

	err = m.db.DeleteAllReasons(height, subspaceID)
	if err != nil {
		return err
	}

	for _, reason := range reasons {
		log.Info().Uint64("subspace", reason.SubspaceID).Uint32("reason", reason.ID).Msg("refreshing reason")

		err = m.db.SaveReason(reason)
		if err != nil {
			return err
		}
	}

	return nil
}

// queryAllReasons queries all the reasons for the given subspace from the node
func (m *Module) queryAllReasons(height int64, subspaceID uint64) ([]types.Reason, error) {
	var reasons []types.Reason

	var nextKey []byte
	var stop = false
	for !stop {
		res, err := m.client.Reasons(
			remote.GetHeightRequestContext(context.Background(), height),
			&reportstypes.QueryReasonsRequest{
				SubspaceId: subspaceID,
				Pagination: &query.PageRequest{
					Key: nextKey,
				},
			},
		)
		if err != nil {
			return nil, err
		}

		for _, reason := range res.Reasons {
			reasons = append(reasons, types.NewReason(reason, height))
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
	}

	return reasons, nil
}
