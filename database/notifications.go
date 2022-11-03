package database

import (
	"encoding/json"

	"github.com/desmos-labs/djuno/v2/types"
)

// SaveNotification stores the given notification inside the database
func (db *Db) SaveNotification(notification types.Notification) error {
	dataBz, err := json.Marshal(&notification.Data)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO notification (user_address, type, data, timestamp) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_user_notification DO UPDATE 
    SET user_address = excluded.user_address,
        data = excluded.data,
        timestamp = excluded.timestamp
WHERE notification.timestamp <= excluded.timestamp`
	_, err = db.SQL.Exec(stmt, notification.RecipientAddress, notification.Type, string(dataBz), notification.Timestamp)
	return err
}
