package notifications

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	juno "github.com/forbole/juno/v4/types"
)

// HandleTx implements modules.TransactionModule
func (m *Module) HandleTx(tx *juno.Tx) error {
	// Get the signers
	signers := map[string]bool{}
	for _, msg := range tx.GetMsgs() {
		if msgExec, isMsgExec := msg.(*authz.MsgExec); isMsgExec {
			for _, msg := range msgExec.Msgs {
				var sdkMsg sdk.Msg
				err := m.cdc.UnpackAny(msg, &sdkMsg)
				if err != nil {
					return err
				}

				// Get all the signers
				signers = getMsgSigners(sdkMsg, signers)
			}
		}

		// Get all the signers
		signers = getMsgSigners(msg, signers)
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

func getMsgSigners(msg sdk.Msg, signers map[string]bool) map[string]bool {
	for _, signer := range msg.GetSigners() {
		if _, ok := signers[signer.String()]; !ok {
			signers[signer.String()] = true
		}
	}
	return signers
}
