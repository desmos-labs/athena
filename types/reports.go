package types

import (
	reportstypes "github.com/desmos-labs/desmos/v4/x/reports/types"
)

type Report struct {
	reportstypes.Report
	Height int64
}

func NewReport(report reportstypes.Report, height int64) Report {
	return Report{
		Report: report,
		Height: height,
	}
}

type Reason struct {
	reportstypes.Reason
	Height int64
}

func NewReason(reason reportstypes.Reason, height int64) Reason {
	return Reason{
		Reason: reason,
		Height: height,
	}
}

type ReportsParams struct {
	reportstypes.Params
	Height int64
}

func NewReportsParams(params reportstypes.Params, height int64) ReportsParams {
	return ReportsParams{
		Params: params,
		Height: height,
	}
}
