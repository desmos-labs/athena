package tips

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/types/query"
	subspacestypes "github.com/desmos-labs/desmos/v4/x/subspaces/types"
	"github.com/forbole/juno/v4/node/remote"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

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
	return m.refreshTips(height, contractAddress)
}

// getContractAddress returns the tips contract address for the given subspace at the provided height
func (m *Module) getContractAddress(height int64, subspaceID uint64) (string, error) {
	// Get all the contracts that match the given code
	var stop = false
	var nextKey []byte
	var contractAddresses []string
	for !stop {
		res, err := m.wasmClient.ContractsByCode(
			remote.GetHeightRequestContext(context.Background(), height),
			&wasmtypes.QueryContractsByCodeRequest{
				CodeId: m.cfg.CodeID,
				Pagination: &query.PageRequest{
					Limit: 100,
					Key:   nextKey,
				},
			},
		)
		if err != nil {
			return "", err
		}

		nextKey = res.Pagination.NextKey
		stop = nextKey == nil
		contractAddresses = append(contractAddresses, res.Contracts...)
	}

	// Search among the contracts addresses, the one for the given subspace
	for _, address := range contractAddresses {
		config, err := m.getContractConfig(height, address)
		if err != nil {
			return "", err
		}

		contractSubspaceID, err := subspacestypes.ParseSubspaceID(config.SubspaceID)
		if err != nil {
			return "", err
		}

		if contractSubspaceID == subspaceID {
			return address, nil
		}
	}

	return "", nil
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
