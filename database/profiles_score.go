package database

import (
	"encoding/json"
	"fmt"

	"github.com/desmos-labs/djuno/v2/types"
)

// SaveApplicationLinkScore stores the given profile score inside the database
func (db *Db) SaveApplicationLinkScore(score *types.ProfileScore) error {
	// Get the id of the application link row associated to this score
	applicationLinkRowID, err := db.getApplicationLinkRowID(score.DesmosAddress, score.Application, score.Username)
	if err != nil {
		return err
	}
	if !applicationLinkRowID.Valid {
		return fmt.Errorf("no application found; address: %s, application: %s, username: %s",
			score.DesmosAddress, score.Application, score.Username)
	}

	// Marshal the details
	detailsBz, err := json.Marshal(&score.Details)
	if err != nil {
		return err
	}

	// Get the score value out of 100
	scoreValue := score.Details.GetScore()
	if scoreValue > 100 {
		scoreValue = 100
	}

	stmt := `
INSERT INTO application_link_score (application_link_row_id, details, score, timestamp)
VALUES ($1, $2, $3, $4)
ON CONFLICT (application_link_row_id) DO UPDATE 
    SET details = excluded.details,
        score = excluded.score,
        timestamp = excluded.timestamp
WHERE application_link_score.timestamp <= excluded.timestamp`
	_, err = db.SQL.Exec(stmt, applicationLinkRowID, string(detailsBz), scoreValue, score.Timestamp)
	return err
}
