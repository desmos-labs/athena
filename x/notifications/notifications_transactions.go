package notifications

import (
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
)

// SendTransactionNotifications notifies the user involved in the transaction
func (m *Module) SendTransactionNotifications(tx *juno.Tx, user string) error {
	data := map[string]string{
		NotificationTypeKey: TypeTransactionSuccess,
		TransactionHashKey:  tx.TxHash,
	}

	if !tx.Successful() {
		data[NotificationTypeKey] = TypeTransactionFailed
		data[TransactionErrorKey] = tx.RawLog
	}

	// Send a notification to the original post owner
	log.Debug().Str("recipient", user).Str("tx_hash", tx.TxHash).Msg("sending notification")
	return m.sendNotification(user, nil, data)
}
