package database

import (
	"encoding/json"
	"fmt"

	"github.com/lib/pq"

	"github.com/desmos-labs/athena/v2/types"
)

// SaveReport saves the given report data inside the database
func (db *Db) SaveReport(report types.Report) error {
	stmt := `
INSERT INTO report (subspace_id, id, message, reporter_address, target, creation_date, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT ON CONSTRAINT unique_subspace_report DO UPDATE 
    SET message = excluded.message,
        reporter_address = excluded.reporter_address,
        target = excluded.target,
        creation_date = excluded.creation_date,
        height = excluded.height
WHERE report.height <= excluded.height
RETURNING row_id`

	targetBz, err := db.cdc.MarshalJSON(report.Target)
	if err != nil {
		return fmt.Errorf("failed to json encode report target: %s", err)
	}

	var reportRowID uint64
	err = db.SQL.QueryRow(stmt,
		report.SubspaceID,
		report.ID,
		report.Message,
		report.Reporter,
		string(targetBz),
		report.CreationDate,
		report.Height,
	).Scan(&reportRowID)
	if err != nil {
		return err
	}

	err = db.insertReportReasons(reportRowID, report.SubspaceID, report.ReasonsIDs)
	if err != nil {
		return err
	}

	return err
}

func (db *Db) insertReportReasons(reportRowID uint64, subspaceID uint64, reasonsIDs []uint32) error {
	if reasonsIDs == nil {
		return nil
	}

	reasonRowsIDs := make(pq.Int64Array, len(reasonsIDs))
	for i, reasonID := range reasonsIDs {
		rowID, err := db.getReasonRowID(subspaceID, reasonID)
		if err != nil {
			return err
		}
		reasonRowsIDs[i] = rowID
	}

	stmt := `INSERT INTO report_reason (report_row_id, reason_row_id) VALUES `

	var vars []interface{}
	for i, reasonRowID := range reasonRowsIDs {
		ei := i * 2
		stmt += fmt.Sprintf(`($%d, $%d),`, ei+1, ei+2)
		vars = append(vars, reportRowID, reasonRowID)
	}

	stmt = stmt[:len(stmt)-1] // Trim trailing ,
	stmt += `ON CONFLICT DO NOTHING`

	_, err := db.SQL.Exec(stmt, vars...)
	return err
}

// DeleteReport removes the report with the given id from the database
func (db *Db) DeleteReport(height int64, subspaceID uint64, reportID uint64) error {
	stmt := `DELETE FROM report WHERE subspace_id = $1 AND id = $2 AND height <= $3`
	_, err := db.SQL.Exec(stmt, subspaceID, reportID, height)
	return err
}

// DeleteAllReports removes all the reports from the database
func (db *Db) DeleteAllReports(height int64, subspaceID uint64) error {
	stmt := `DELETE FROM report WHERE subspace_id = $1 AND height <= $2`
	_, err := db.SQL.Exec(stmt, subspaceID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

func (db *Db) getReasonRowID(subspaceID uint64, reasonID uint32) (int64, error) {
	stmt := `SELECT row_id FROM subspace_report_reason WHERE subspace_id = $1 AND id = $2`

	var rowID int64
	err := db.SQL.QueryRow(stmt, subspaceID, reasonID).Scan(&rowID)
	if err != nil {
		return 0, err
	}

	return rowID, nil
}

// SaveReason saves the given reason insinde the database
func (db *Db) SaveReason(reason types.Reason) error {
	stmt := `
INSERT INTO subspace_report_reason (subspace_id, id, title, description, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT ON CONSTRAINT unique_subspace_reason DO UPDATE 
    SET title = excluded.title,
        description = excluded.description,
        height = excluded.height
WHERE subspace_report_reason.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, reason.SubspaceID, reason.ID, reason.Title, reason.Description, reason.Height)
	return err
}

// DeleteReason removes the reason having the given id from the database along with all the associated reports
func (db *Db) DeleteReason(height int64, subspaceID uint64, reasonID uint32) error {
	// Delete the reason
	stmt := `DELETE FROM subspace_report_reason WHERE subspace_id = $1 AND id = $2 AND height <= $2`
	_, err := db.SQL.Exec(stmt, subspaceID, reasonID, height)
	return err
}

// DeleteAllReasons deletes all the reasons from the database
func (db *Db) DeleteAllReasons(height int64, subspaceID uint64) error {
	stmt := `DELETE FROM subspace_report_reason WHERE subspace_id = $1 AND height <= $2`
	_, err := db.SQL.Exec(stmt, subspaceID, height)
	return err
}

// --------------------------------------------------------------------------------------------------------------------

// SaveReportsParams saves the given reports params inside the database
func (db *Db) SaveReportsParams(params types.ReportsParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling reports params: %s", err)
	}

	stmt := `
INSERT INTO reports_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE reports_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing reports params: %s", err)
	}

	return nil
}
