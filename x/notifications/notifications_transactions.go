package notifications

import (
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/v2/types"
)

// SendTransactionNotifications notifies the user involved in the transaction
func (m *Module) SendTransactionNotifications(tx *juno.Tx, user string) error {
	data := map[string]string{
		types.NotificationTypeKey: types.TypeTransactionSuccess,
		types.TransactionHashKey:  tx.TxHash,
	}

	if !tx.Successful() {
		data[types.NotificationTypeKey] = types.TypeTransactionFailed
		data[types.TransactionErrorKey] = tx.RawLog
	}

	// Send a notification to the original post owner
	log.Trace().Str("module", m.Name()).Str("recipient", user).Str("tx hash", tx.TxHash).
		Str("notification type", "transaction").Msg("sending notification")

	return m.SendAndStoreNotification(types.NewNotificationUserRecipient(user), types.NewStdNotificationDataWithConfig(nil, data))
}
