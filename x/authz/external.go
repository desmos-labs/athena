package authz

import (
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"

	"github.com/desmos-labs/djuno/v2/utils"
)

// RefreshAuthorizations refreshes all the authorizations data
func (m *Module) RefreshAuthorizations(height int64) error {
	return m.refreshAuthorizations(height)
}

// refreshAuthorizations queries all the MsgGrant and MsgRevoke transactions to refresh the stored authorizations
func (m *Module) refreshAuthorizations(height int64) error {
	// Get all the grant transactions
	grantsQuery := fmt.Sprintf("%s.%s CONTAINS '%s' AND tx.height <= %d",
		proto.MessageName(&authz.EventGrant{}),
		"granter",
		"desmos",
		height,
	)
	txs, err := utils.QueryTxs(m.node, grantsQuery)
	if err != nil {
		return err
	}

	// Get all the revoke transactions
	revokeQuery := fmt.Sprintf("%s.%s CONTAINS '%s' AND tx.height <= %d",
		proto.MessageName(&authz.EventRevoke{}),
		"granter",
		"desmos",
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
