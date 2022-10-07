package tips

import (
	"encoding/hex"
	"fmt"
	"sort"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"

	"github.com/desmos-labs/djuno/v2/types"
	"github.com/desmos-labs/djuno/v2/utils"
)

// RefreshData refreshes the data related to the tips contract for the given subspace, if any
func (m *Module) RefreshData(height int64, subspaceID uint64) error {
	contractAddress, err := m.getContractAddress(height, subspaceID)
	if err != nil {
		return err
	}

	// Make sure there is a contract for this subspace
	if contractAddress == "" {
		return nil
	}

	// Store the contract
	err = m.db.SaveContract(types.NewContract(contractAddress, types.ContractTypeTips, height))
	if err != nil {
		return err
	}

	// Refresh the tips data
	return m.refreshTips(height, contractAddress)
}

// getContractAddress returns the tips contract address for the given subspace at the provided height
func (m *Module) getContractAddress(height int64, subspaceID uint64) (string, error) {
	// Query all the transactions
	permissionsQuery := fmt.Sprintf("%s.%s=%d AND %s.%s='%d' AND tx.height <= %d",
		wasmtypes.EventTypeInstantiate,
		wasmtypes.AttributeKeyCodeID,
		m.cfg.CodeID,
		wasmtypes.WasmModuleEventType,
		subspacestypes.AttributeKeySubspaceID,
		subspaceID,
		height,
	)
	txs, err := utils.QueryTxs(m.node, permissionsQuery)
	if err != nil {
		return "", err
	}

	// If there are not transactions, just return
	if len(txs) == 0 {
		return "", nil
	}

	// Sort the txs based on their descending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height > txs[j].Height
	})

	// Get the transaction details
	transaction, err := m.node.Tx(hex.EncodeToString(txs[0].Tx.Hash()))
	if err != nil {
		return "", err
	}

	// Handle only the MsgInstantiateContract
	for index, msg := range transaction.GetMsgs() {
		if _, ok := msg.(*wasmtypes.MsgInstantiateContract); ok {
			return m.base.ParseContractAddress(transaction, index)
		}
	}

	return "", fmt.Errorf("no MsgInstantiateContract found")
}

// refreshTips fetches and stores all the tips sent using the contract having the given address
// before or on the provided height
func (m *Module) refreshTips(height int64, contractAddress string) error {
	// Query all the transactions
	permissionsQuery := fmt.Sprintf("%s.%s='%s' AND %s.%s='%s' AND tx.height <= %d",
		wasmtypes.WasmModuleEventType,
		sdk.AttributeKeyAction,
		"send_tip",
		wasmtypes.WasmModuleEventType,
		wasmtypes.AttributeKeyContractAddr,
		contractAddress,
		height,
	)
	txs, err := utils.QueryTxs(m.node, permissionsQuery)
	if err != nil {
		return err
	}

	// Sort the txs based on their ascending height
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Height < txs[j].Height
	})

	for _, tx := range txs {
		transaction, err := m.node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle only the MsgSetUserPermissions
		for index, msg := range transaction.GetMsgs() {
			if msg, ok := msg.(*authz.MsgExec); ok {
				innerMsgs, err := msg.GetMessages()
				if err != nil {
					return err
				}

				for innerIndex, innerMsg := range innerMsgs {
					err = m.HandleMsgExec(index, msg, innerIndex, innerMsg, transaction)
					if err != nil {
						return err
					}
				}
			}

			if msg, ok := msg.(*wasmtypes.MsgExecuteContract); ok {
				err = m.handleMsgExecuteContract(transaction, msg)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
