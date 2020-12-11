package notifications

import (
	juno "github.com/desmos-labs/juno/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/abci/types"
)

const (
	TypeTransactionSuccess = "transaction_success"
	TypeTransactionFailed  = "transaction_fail"

	TransactionHashKey  = "tx_hash"
	TransactionErrorKey = "tx_error"
)

// SendTransactionResultNotification sends to the given user a notification telling him
// that the specified transaction has either succeeded or failed
func SendTransactionResultNotification(tx *juno.Tx, user string) error {
	result := TypeTransactionSuccess
	if tx.Code != tmtypes.CodeTypeOK {
		result = TypeTransactionFailed
	}

	log.Info().
		Str("recipient", user).
		Str("tx_hash", tx.TxHash).
		Str("tx_result", result).
		Msg("sending notification")

	data := map[string]string{
		NotificationTypeKey: result,

		TransactionHashKey:  tx.TxHash,
		TransactionErrorKey: tx.RawLog,
	}

	// Send a notification to the original post owner
	return SendNotification(user, nil, data)
}
