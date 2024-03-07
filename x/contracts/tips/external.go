package tips

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	subspacestypes "github.com/desmos-labs/desmos/v7/x/subspaces/types"

	"github.com/desmos-labs/athena/v2/types"
	"github.com/desmos-labs/athena/v2/utils"
)

// RefreshData refreshes the data related to the tips contract for the given subspace, if any
func (m *Module) RefreshData(height int64, subspaceID uint64) error {
	for _, contractAddress := range m.cfg.Addresses {
		isSubspaceContract, err := m.isSubspaceContract(height, subspaceID, contractAddress)
		if err != nil {
			return err
		}

		if !isSubspaceContract {
			continue
		}

		// Get the contract config
		config, err := m.getContractConfig(height, contractAddress)
		if err != nil {
			return err
		}

		configBz, err := json.Marshal(&config)
		if err != nil {
			return err
		}

		// Store the contract
		err = m.db.SaveContract(types.NewContract(contractAddress, types.ContractTypeTips, configBz, height))
		if err != nil {
			return err
		}

		// Refresh the tips data
		err = m.refreshTips(height, contractAddress)
		if err != nil {
			return err
		}
	}

	return nil
}

// isSubspaceContract tells whether the contract having the given address is related to the provided subspace
func (m *Module) isSubspaceContract(height int64, subspaceID uint64, contractAddress string) (bool, error) {
	config, err := m.getContractConfig(height, contractAddress)
	if err != nil {
		return false, err
	}

	contractSubspaceID, err := subspacestypes.ParseSubspaceID(config.SubspaceID)
	if err != nil {
		return false, err
	}

	return contractSubspaceID == subspaceID, nil
}

// refreshContractConfig refreshes the configuration for the contract having the given address at the provided height
func (m *Module) refreshContractConfig(height int64, address string) error {
	config, err := m.getContractConfig(height, address)
	if err != nil {
		return err
	}

	configBz, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return m.db.SaveContract(types.NewContract(address, types.ContractTypeTips, configBz, height))
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
