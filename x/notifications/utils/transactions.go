package utils

import (
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)

const (
	NotificationTypeTxSuccess = "transaction_success"
	NotificationTypeTxFailed  = "transaction_fail"

	AttributeKeyTxHash  = "tx_hash"
	AttributeKeyTxError = "tx_error"
)

// SendTransactionResultNotification sends to the given user a notification telling him
// that the specified transaction has either succeeded or failed
func SendTransactionResultNotification(tx *juno.Tx, user string) error {
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
	return SendNotification(user, nil, data)
}
