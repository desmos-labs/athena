package notifications

import (
	juno "github.com/forbole/juno/v3/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
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
		err := m.SendTransactionNotifications(tx, signer)
		if err != nil {
			return err
		}
	}

	return nil
}
