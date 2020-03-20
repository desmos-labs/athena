package notifications

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/types"
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
func SendTransactionResultNotification(tx types.Tx, user sdk.AccAddress) error {
	result := TypeTransactionSuccess
	if tx.Code != tmtypes.CodeTypeOK {
		result = TypeTransactionFailed
	}

	log.Info().Msg(fmt.Sprintf("Sending failed tx notification to %s for tx with hash %s, result: %s",
		user, tx.TxHash, result))

	data := map[string]string{
		NotificationTypeKey: result,

		TransactionHashKey:  tx.TxHash,
		TransactionErrorKey: tx.RawLog,
	}

	// Send a notification to the original post owner
	return SendNotification(user.String(), nil, data)
}
