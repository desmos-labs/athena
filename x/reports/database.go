package reports

import "github.com/desmos-labs/athena/types"

type Database interface {
	SaveReport(report types.Report) error
	DeleteReport(height int64, subspaceID uint64, reportID uint64) error
	DeleteAllReports(height int64, subspaceID uint64) error
	SaveReason(reason types.Reason) error
	DeleteReason(height int64, subspaceID uint64, reasonID uint32) error
	DeleteAllReasons(height int64, subspaceID uint64) error
	SaveReportsParams(params types.ReportsParams) error
}
