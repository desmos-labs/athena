package notifications

import (
	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/djuno/v2/x/notifications/builder"
)

// SendTransactionNotifications notifies the user involved in the transaction
func (m *Module) SendTransactionNotifications(tx *juno.Tx, user string) error {
	data := map[string]string{
		builder.NotificationTypeKey: builder.TypeTransactionSuccess,
		builder.TransactionHashKey:  tx.TxHash,
	}

	if !tx.Successful() {
		data[builder.NotificationTypeKey] = builder.TypeTransactionFailed
		data[builder.TransactionErrorKey] = tx.RawLog
	}

	// Send a notification to the original post owner
	log.Debug().Str("module", m.Name()).Str("recipient", user).Str("tx hash", tx.TxHash).
		Str("notification type", data[builder.NotificationTypeKey]).Msg("sending notification")
	return m.SendNotification(user, nil, data)
}
