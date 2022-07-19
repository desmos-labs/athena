package notifications

import (
	juno "github.com/forbole/juno/v3/types"
	"github.com/rs/zerolog/log"
)

// HandleTx implements modules.TransactionModule
func (m Module) HandleTx(tx *juno.Tx) error {
	// Get the signers
	var signers map[string]bool
	for _, msg := range tx.GetMsgs() {
		for _, signer := range msg.GetSigners() {
			// Add the signer to the map
			if _, ok := signers[signer.String()]; !ok {
				signers[signer.String()] = true
			}
		}
	}

	// Send the transaction result notification to all the signers
	for signer := range signers {
		err := m.sendTransactionResultNotification(tx, signer)
		if err != nil {
			return err
		}
	}

	return nil
}

// sendTransactionResultNotification notifies the user involved in the transaction
func (m *Module) sendTransactionResultNotification(tx *juno.Tx, user string) error {
	result := TypeTransactionSuccess
	if !tx.Successful() {
		result = TypeTransactionFailed
	}

	data := map[string]string{
		NotificationTypeKey: result,
		TransactionHashKey:  tx.TxHash,
		TransactionErrorKey: tx.RawLog,
	}

	// Send a notification to the original post owner
	log.Debug().Str("recipient", user).Str("tx_hash", tx.TxHash).Msg("sending notification")
	return m.sendNotification(user, nil, data)
}
