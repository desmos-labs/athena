package notifications

import (
	juno "github.com/forbole/juno/v2/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)

const (
	NotificationTypeTxSuccess = "transaction_success"
	NotificationTypeTxFailed  = "transaction_fail"

	AttributeKeyTxHash  = "tx_hash"
	AttributeKeyTxError = "tx_error"
)

// sendTransactionResultNotification sends to the given user a notification telling him
// that the specified transaction has either succeeded or failed
func (m *Module) sendTransactionResultNotification(tx *juno.Tx, user string) error {
	result := NotificationTypeTxSuccess
	if tx.Code != tmtypes.CodeTypeOK {
		result = NotificationTypeTxFailed
	}

	log.Info().
		Str("recipient", user).
		Str("tx_hash", tx.TxHash).
		Str("tx_result", result).
		Msg("sending notification")

	data := map[string]string{
		NotificationTypeKey: result,

		AttributeKeyTxHash:  tx.TxHash,
		AttributeKeyTxError: tx.RawLog,
	}

	// Send a notification to the original post owner
	return m.sendNotification(user, nil, data)
}
