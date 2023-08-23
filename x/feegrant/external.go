package feegrant

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/desmos-labs/desmos/v6/app"

	"github.com/desmos-labs/djuno/v2/utils"
)

// RefreshFeeGrants refreshes all the fee grants
func (m *Module) RefreshFeeGrants(height int64) error {
	return m.refreshFeeGrants(height)
}

// refreshFeeGrants searches all the MsgGrantAllowance and MsgRevokeAllowance messages and stores the related data
func (m *Module) refreshFeeGrants(height int64) error {
	// Get all the grant transactions
	grantsQuery := fmt.Sprintf("%s.%s CONTAINS '%s' AND tx.height <= %d",
		feegrant.EventTypeSetFeeGrant,
		feegrant.AttributeKeyGranter,
		app.Bech32MainPrefix,
		height,
	)
	txs, err := utils.QueryTxs(m.node, grantsQuery)
	if err != nil {
		return err
	}

	// Get all the revoke transactions
	revokeQuery := fmt.Sprintf("%s.%s CONTAINS '%s' AND tx.height <= %d",
		feegrant.EventTypeRevokeFeeGrant,
		feegrant.AttributeKeyGranter,
		app.Bech32MainPrefix,
		height,
	)
	revokeTxs, err := utils.QueryTxs(m.node, revokeQuery)
	if err != nil {
		return err
	}

	// combine the transactions
	txs = append(txs, revokeTxs...)

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	// Handle all the transactions
	for _, tx := range txs {
		transaction, err := m.node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		for index, msg := range transaction.GetMsgs() {
			err = m.HandleMsg(index, msg, transaction)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
