package database

import (
	reportstypes "github.com/desmos-labs/desmos/x/staging/reports/types"
)

// SaveReport allows to store the given report properly
func (db DesmosDb) SaveReport(report reportstypes.Report) error {
	err := db.SaveUserIfNotExisting(report.User)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO report(post_id, type, message, reporter_address) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err = db.Sql.Exec(stmt, report.PostId, report.Type, report.Message, report.User)
	return err
}
