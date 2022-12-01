package database

import (
	"encoding/json"
	"time"

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

// SaveToken stores the given notification token inside the database
func (db *Db) SaveToken(token types.NotificationToken) error {
	stmt := `
INSERT INTO notification_token (user_address, device_token, timestamp)
VALUES ($1, $2, $3)`
	_, err := db.SQL.Exec(stmt, token.UserAddress, token.Token, token.Timestamp)
	return err
}

type notificationTokenRow struct {
	UserAddress string    `db:"user_address"`
	DeviceToken string    `db:"device_token"`
	Timestamp   time.Time `db:"timestamp"`
}

// GetUserTokens returns all the notifications tokens associated to all the devices of the user having the given address
func (db *Db) GetUserTokens(userAddress string) ([]types.NotificationToken, error) {
	stmt := `SELECT * FROM notification_token WHERE user_address = $1`

	var rows []notificationTokenRow
	err := db.SQL.Select(&rows, stmt, userAddress)
	if err != nil {
		return nil, err
	}

	tokens := make([]types.NotificationToken, len(rows))
	for i, row := range rows {
		tokens[i] = types.NewNotificationToken(row.UserAddress, row.DeviceToken, row.Timestamp)
	}
	return tokens, nil
}
